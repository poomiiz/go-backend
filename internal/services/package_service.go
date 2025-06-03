package services

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/poomiiz/go-backend/internal/utils"
)

// Package โครงสร้างข้อมูลของแต่ละ Package ใน Firestore
type Package struct {
	Name         string    `firestore:"name"`
	CoinCost     int64     `firestore:"coinCost"`
	DurationDays int       `firestore:"durationDays"`
	CreatedAt    time.Time `firestore:"createdAt"`
	UpdatedAt    time.Time `firestore:"updatedAt"`
}

// UserPackage โครงสร้างข้อมูลการใช้งาน Package ของผู้ใช้
type UserPackage struct {
	UserID    string    `firestore:"userId"`
	PackageID string    `firestore:"packageId"`
	StartedAt time.Time `firestore:"startedAt"`
	ExpiresAt time.Time `firestore:"expiresAt"`
	CreatedAt time.Time `firestore:"createdAt"`
	UpdatedAt time.Time `firestore:"updatedAt"`
}

type PackageService struct {
	pkgCol     *firestore.CollectionRef
	userPkgCol *firestore.CollectionRef
	coinSvc    *CoinService
}

func NewPackageService(coinSvc *CoinService) *PackageService {
	return &PackageService{
		pkgCol:     utils.Client.Collection("packages"),
		userPkgCol: utils.Client.Collection("user_packages"),
		coinSvc:    coinSvc,
	}
}

// CreatePackage: สำหรับ Admin สร้าง Package ใหม่
func (s *PackageService) CreatePackage(ctx context.Context, name string, coinCost int64, durationDays int) (string, error) {
	now := time.Now()
	pkg := Package{
		Name:         name,
		CoinCost:     coinCost,
		DurationDays: durationDays,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	docRef, _, err := s.pkgCol.Add(ctx, pkg)
	if err != nil {
		return "", err
	}
	return docRef.ID, nil
}

// BuyPackage: ผู้ใช้ซื้อหรือต่ออายุ Package
func (s *PackageService) BuyPackage(ctx context.Context, userID, packageID string) (*UserPackage, error) {
	// 1. ดึงข้อมูล Package ออกมาจาก Firestore
	pkgSnap, err := s.pkgCol.Doc(packageID).Get(ctx)
	if err != nil {
		return nil, err
	}
	var pkg Package
	if err := pkgSnap.DataTo(&pkg); err != nil {
		return nil, err
	}

	// 2. ตรวจสอบยอดเหรียญ (Deduct ลบยอดก่อน)
	if err := s.coinSvc.Deduct(ctx, userID, pkg.CoinCost); err != nil {
		return nil, err
	}

	// 3. หา UserPackage เดิม (ถ้ามี) เพื่อดูว่ายัง active หรือหมดอายุ
	q := s.userPkgCol.Where("userId", "==", userID).Where("packageId", "==", packageID).Limit(1)
	docs, err := q.Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	var upsnap *firestore.DocumentSnapshot
	if len(docs) > 0 {
		upsnap = docs[0]
	}

	// คำนวณเวลาหมดอายุใหม่
	var newExpiry time.Time
	if upsnap != nil {
		var existing UserPackage
		if err := upsnap.DataTo(&existing); err != nil {
			return nil, err
		}
		if existing.ExpiresAt.After(now) {
			// ยังไม่หมดอายุ ให้ต่อจากวันเดิม
			newExpiry = existing.ExpiresAt.AddDate(0, 0, pkg.DurationDays)
		} else {
			// หมดอายุไปแล้ว ให้เริ่มใหม่จาก now
			newExpiry = now.AddDate(0, 0, pkg.DurationDays)
		}
	} else {
		// ไม่เคยซื้อมาก่อน
		newExpiry = now.AddDate(0, 0, pkg.DurationDays)
	}

	if upsnap == nil {
		// สร้าง document ใหม
		userPkg := UserPackage{
			UserID:    userID,
			PackageID: packageID,
			StartedAt: now,
			ExpiresAt: newExpiry,
			CreatedAt: now,
			UpdatedAt: now,
		}
		docRef, _, err := s.userPkgCol.Add(ctx, userPkg)
		if err != nil {
			return nil, err
		}
		userPkgDoc, err := docRef.Get(ctx)
		if err != nil {
			return nil, err
		}
		var created UserPackage
		if err := userPkgDoc.DataTo(&created); err != nil {
			return nil, err
		}
		return &created, nil
	}

	// อัปเดต document เดิม
	updates := []firestore.Update{
		{Path: "startedAt", Value: now},
		{Path: "expiresAt", Value: newExpiry},
		{Path: "updatedAt", Value: now},
	}
	_, err = upsnap.Ref.Update(ctx, updates)
	if err != nil {
		return nil, err
	}
	var updated UserPackage
	snapAfter, err := upsnap.Ref.Get(ctx)
	if err != nil {
		return nil, err
	}
	if err := snapAfter.DataTo(&updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

// CheckUserPackage: ตรวจสอบว่าผู้ใช้ยังมี Package ไหน active อยู่หรือไม่
func (s *PackageService) CheckUserPackage(ctx context.Context, userID string) (bool, error) {
	now := time.Now()
	q := s.userPkgCol.Where("userId", "==", userID).Where("expiresAt", ">", now).Limit(1)
	docs, err := q.Documents(ctx).GetAll()
	if err != nil {
		return false, err
	}
	return len(docs) > 0, nil
}

// GetUserPackages: ดึงข้อมูล UserPackage ของผู้ใช้ทั้งหมด
func (s *PackageService) GetUserPackages(ctx context.Context, userID string) ([]UserPackage, error) {
	q := s.userPkgCol.Where("userId", "==", userID)
	docs, err := q.Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}
	var result []UserPackage
	for _, doc := range docs {
		var up UserPackage
		if err := doc.DataTo(&up); err != nil {
			return nil, err
		}
		result = append(result, up)
	}
	return result, nil
}
