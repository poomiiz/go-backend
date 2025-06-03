package services

import (
	"context"
	"errors"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/poomiiz/go-backend/internal/utils"
)

// Payment โครงสร้างข้อมูลใน Firestore
type Payment struct {
	UserID        string    `firestore:"userId"`
	Amount        int64     `firestore:"amount"`
	Provider      string    `firestore:"provider"`      // e.g. "Omise", "TrueMoney"
	ProviderRefID string    `firestore:"providerRefId"` // reference จากผู้ให้บริการ
	Status        string    `firestore:"status"`        // "pending", "paid", "failed"
	Commission    float64   `firestore:"commission"`
	CreatedAt     time.Time `firestore:"createdAt"`
	UpdatedAt     time.Time `firestore:"updatedAt"`
}

type PaymentService struct {
	col               *firestore.CollectionRef
	commissionPercent float64
}

func NewPaymentService(commissionPercent float64) *PaymentService {
	return &PaymentService{
		col:               utils.Client.Collection("payments"),
		commissionPercent: commissionPercent,
	}
}

// CreatePayment: สร้าง Payment ใหม่ (status = "pending")
func (s *PaymentService) CreatePayment(ctx context.Context, userID string, amount int64, provider string, providerRefID string) (string, error) {
	now := time.Now()
	pay := Payment{
		UserID:        userID,
		Amount:        amount,
		Provider:      provider,
		ProviderRefID: providerRefID,
		Status:        "pending",
		Commission:    0.0,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	docRef, _, err := s.col.Add(ctx, pay)
	if err != nil {
		return "", err
	}
	return docRef.ID, nil
}

// VerifyPayment: ตรวจสอบสถานะกับ Provider แล้วอัปเดตใน Firestore
func (s *PaymentService) VerifyPayment(ctx context.Context, paymentID string) error {
	docSnap, err := s.col.Doc(paymentID).Get(ctx)
	if err != nil {
		return err
	}
	var p Payment
	if err := docSnap.DataTo(&p); err != nil {
		return err
	}
	if p.Status != "pending" {
		return errors.New("payment is not in pending state")
	}

	// 1. ตรวจสอบกับ Provider จริง (ตัวอย่างเป็น pseudo-code)
	paid, err := s.checkProvider(p.Provider, p.ProviderRefID)
	if err != nil {
		return err
	}

	// 2. คำนวณ Commission ถ้า paid
	newStatus := "failed"
	if paid {
		newStatus = "paid"
	}
	commission := 0.0
	if paid {
		commission = float64(p.Amount) * s.commissionPercent / 100.0
	}

	// 3. อัปเดต document
	updates := []firestore.Update{
		{Path: "status", Value: newStatus},
		{Path: "commission", Value: commission},
		{Path: "updatedAt", Value: time.Now()},
	}
	_, err = docSnap.Ref.Update(ctx, updates)
	return err
}

// checkProvider เป็นตัวอย่าง pseudo-code สำหรับตรวจสอบกับ Provider
func (s *PaymentService) checkProvider(provider, refID string) (bool, error) {
	// TODO: เขียนโค้ดเช็คจาก Omise หรือ TrueMoney
	// เช่น HTTP request ไปหา API ของ Provider แล้ว parse status
	// ถ้า status == "successful" => return true
	return false, nil
}
