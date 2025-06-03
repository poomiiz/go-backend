package services

import (
	"context"
	"errors"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/poomiiz/go-backend/internal/utils"
)

type CoinBalance struct {
	UserID    string    `firestore:"userId"`
	Balance   int64     `firestore:"balance"`
	UpdatedAt time.Time `firestore:"updatedAt"`
}

type CoinService struct {
	col *firestore.CollectionRef
}

func NewCoinService() *CoinService {
	return &CoinService{
		col: utils.Client.Collection("coin_balances"),
	}
}

// GetBalance ดึงยอดเหรียญของผู้ใช้
func (s *CoinService) GetBalance(ctx context.Context, userID string) (int64, error) {
	q := s.col.Where("userId", "==", userID).Limit(1)
	docs, err := q.Documents(ctx).GetAll()
	if err != nil {
		return 0, err
	}
	if len(docs) == 0 {
		return 0, nil
	}
	var cb CoinBalance
	if err := docs[0].DataTo(&cb); err != nil {
		return 0, err
	}
	return cb.Balance, nil
}

// TopUp เติมเหรียญให้ผู้ใช้
func (s *CoinService) TopUp(ctx context.Context, userID string, amount int64) error {
	now := time.Now()
	q := s.col.Where("userId", "==", userID).Limit(1)
	docs, err := q.Documents(ctx).GetAll()
	if err != nil {
		return err
	}
	if len(docs) == 0 {
		cb := CoinBalance{UserID: userID, Balance: amount, UpdatedAt: now}
		_, _, err = s.col.Add(ctx, cb)
		return err
	}
	docRef := docs[0].Ref
	_, err = docRef.Update(ctx, []firestore.Update{
		{Path: "balance", Value: firestore.Increment(amount)},
		{Path: "updatedAt", Value: now},
	})
	return err
}

// Deduct หักยอดเหรียญของผู้ใช้
func (s *CoinService) Deduct(ctx context.Context, userID string, amount int64) error {
	now := time.Now()
	q := s.col.Where("userId", "==", userID).Limit(1)
	docs, err := q.Documents(ctx).GetAll()
	if err != nil {
		return err
	}
	if len(docs) == 0 {
		return errors.New("no balance record")
	}
	var cb CoinBalance
	doc := docs[0]
	if err := doc.DataTo(&cb); err != nil {
		return err
	}
	if cb.Balance < amount {
		return errors.New("insufficient balance")
	}
	_, err = doc.Ref.Update(ctx, []firestore.Update{
		{Path: "balance", Value: firestore.Increment(-amount)},
		{Path: "updatedAt", Value: now},
	})
	return err
}

// Transfer โอนเหรียญจากผู้ใช้หนึ่งไปยังอีกคน
func (s *CoinService) Transfer(ctx context.Context, fromUserID, toUserID string, amount int64) error {
	if fromUserID == toUserID {
		return errors.New("cannot transfer to self")
	}
	// หักยอดจากผู้ส่ง
	if err := s.Deduct(ctx, fromUserID, amount); err != nil {
		return err
	}
	// เติมยอดให้ผู้รับ
	return s.TopUp(ctx, toUserID, amount)
}
