package services

import (
	"context"
	"errors"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/poomiiz/go-backend/internal/utils"
)

type Review struct {
	UserID    string    `firestore:"userId"`
	SeerID    string    `firestore:"seerId"`
	Rating    int       `firestore:"rating"`
	Content   string    `firestore:"content"`
	Status    string    `firestore:"status"`
	CreatedAt time.Time `firestore:"createdAt"`
	UpdatedAt time.Time `firestore:"updatedAt"`
}

type Appeal struct {
	ReviewID  string    `firestore:"reviewId"`
	UserID    string    `firestore:"userId"`
	Reason    string    `firestore:"reason"`
	CreatedAt time.Time `firestore:"createdAt"`
	UpdatedAt time.Time `firestore:"updatedAt"`
}

type ReviewService struct {
	reviewCol *firestore.CollectionRef
	appealCol *firestore.CollectionRef
}

// NewReviewService สร้าง instance ของ ReviewService
func NewReviewService() *ReviewService {
	return &ReviewService{
		reviewCol: utils.Client.Collection("reviews"),
		appealCol: utils.Client.Collection("appeals"),
	}
}

// SubmitReview: ผู้ใช้เขียนรีวิวใหม่ (status = "pending")
func (s *ReviewService) SubmitReview(ctx context.Context, userID, seerID string, rating int, content string) (string, error) {
	now := time.Now()
	rev := Review{
		UserID:    userID,
		SeerID:    seerID,
		Rating:    rating,
		Content:   content,
		Status:    "pending",
		CreatedAt: now,
		UpdatedAt: now,
	}
	docRef, _, err := s.reviewCol.Add(ctx, rev)
	if err != nil {
		return "", err
	}
	return docRef.ID, nil
}

// GetPendingReviews: ดึงรีวิวที่ status == "pending"
func (s *ReviewService) GetPendingReviews(ctx context.Context) ([]Review, error) {
	q := s.reviewCol.Where("status", "==", "pending")
	docs, err := q.Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}
	var list []Review
	for _, doc := range docs {
		var r Review
		if err := doc.DataTo(&r); err != nil {
			continue
		}
		list = append(list, r)
	}
	return list, nil
}

// ApproveReview: Admin อนุมัติรีวิว
func (s *ReviewService) ApproveReview(ctx context.Context, reviewID string) error {
	docRef := s.reviewCol.Doc(reviewID)
	snap, err := docRef.Get(ctx)
	if err != nil {
		return err
	}
	var r Review
	if err := snap.DataTo(&r); err != nil {
		return err
	}
	if r.Status != "pending" {
		return errors.New("review not pending")
	}
	_, err = docRef.Update(ctx, []firestore.Update{
		{Path: "status", Value: "approved"},
		{Path: "updatedAt", Value: time.Now()},
	})
	return err
}

// RejectReview: Admin ปฏิเสธรีวิว
func (s *ReviewService) RejectReview(ctx context.Context, reviewID string) error {
	docRef := s.reviewCol.Doc(reviewID)
	snap, err := docRef.Get(ctx)
	if err != nil {
		return err
	}
	var r Review
	if err := snap.DataTo(&r); err != nil {
		return err
	}
	if r.Status != "pending" {
		return errors.New("review not pending")
	}
	_, err = docRef.Update(ctx, []firestore.Update{
		{Path: "status", Value: "rejected"},
		{Path: "updatedAt", Value: time.Now()},
	})
	return err
}

// DeleteReview: Admin ลบรีวิว (ถ้าต้องการ)
func (s *ReviewService) DeleteReview(ctx context.Context, reviewID string) error {
	_, err := s.reviewCol.Doc(reviewID).Delete(ctx)
	return err
}

// AppealReview: ผู้ใช้ยื่นอุทธรณ์รีวิวที่ถูกปฏิเสธ
func (s *ReviewService) AppealReview(ctx context.Context, reviewID, userID, reason string) (string, error) {
	docSnap, err := s.reviewCol.Doc(reviewID).Get(ctx)
	if err != nil {
		return "", err
	}
	var r Review
	if err := docSnap.DataTo(&r); err != nil {
		return "", err
	}
	if r.UserID != userID || r.Status != "rejected" {
		return "", errors.New("cannot appeal this review")
	}
	now := time.Now()
	app := Appeal{
		ReviewID:  reviewID,
		UserID:    userID,
		Reason:    reason,
		CreatedAt: now,
		UpdatedAt: now,
	}
	docRef, _, err := s.appealCol.Add(ctx, app)
	if err != nil {
		return "", err
	}
	return docRef.ID, nil
}
