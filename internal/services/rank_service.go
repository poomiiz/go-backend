package services

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/poomiiz/go-backend/internal/utils"
)

type Rank struct {
	SeerID       string    `firestore:"seerId"`
	RatingSum    int       `firestore:"ratingSum"`
	ReviewCount  int       `firestore:"reviewCount"`
	RankPosition int       `firestore:"rankPosition"`
	UpdatedAt    time.Time `firestore:"updatedAt"`
}

type Commission struct {
	SeerID    string    `firestore:"seerId"`
	Amount    float64   `firestore:"amount"`
	Month     string    `firestore:"month"`
	CreatedAt time.Time `firestore:"createdAt"`
}

type Bonus struct {
	SeerID    string    `firestore:"seerId"`
	Amount    float64   `firestore:"amount"`
	Period    string    `firestore:"period"`
	CreatedAt time.Time `firestore:"createdAt"`
}

type RankService struct {
	reviewCol     *firestore.CollectionRef
	rankCol       *firestore.CollectionRef
	commissionCol *firestore.CollectionRef
	bonusCol      *firestore.CollectionRef
	paymentCol    *firestore.CollectionRef
}

func NewRankService() *RankService {
	return &RankService{
		reviewCol:     utils.Client.Collection("reviews"),
		rankCol:       utils.Client.Collection("ranks"),
		commissionCol: utils.Client.Collection("commissions"),
		bonusCol:      utils.Client.Collection("bonuses"),
		paymentCol:    utils.Client.Collection("payments"),
	}
}

func (s *RankService) CalculateRankings(ctx context.Context, since time.Time) error {
	q := s.reviewCol.Where("status", "==", "approved").Where("createdAt", ">=", since)
	docs, err := q.Documents(ctx).GetAll()
	if err != nil {
		return err
	}
	type agg struct {
		ratingSum   int
		reviewCount int
	}
	aggMap := make(map[string]*agg)
	for _, doc := range docs {
		var r Review
		if err := doc.DataTo(&r); err != nil {
			continue
		}
		if _, ok := aggMap[r.SeerID]; !ok {
			aggMap[r.SeerID] = &agg{}
		}
		aggMap[r.SeerID].ratingSum += r.Rating
		aggMap[r.SeerID].reviewCount++
	}
	// Build slice and sort
	type rItem struct {
		SeerID      string
		RatingSum   int
		ReviewCount int
	}
	var items []rItem
	for seerID, a := range aggMap {
		items = append(items, rItem{SeerID: seerID, RatingSum: a.ratingSum, ReviewCount: a.reviewCount})
	}
	// Sort by RatingSum desc
	for i := 0; i < len(items); i++ {
		for j := i + 1; j < len(items); j++ {
			if items[j].RatingSum > items[i].RatingSum {
				items[i], items[j] = items[j], items[i]
			}
		}
	}
	// Delete old rankings
	oldDocs, _ := s.rankCol.Documents(ctx).GetAll()
	for _, rd := range oldDocs {
		rd.Ref.Delete(ctx)
	}
	now := time.Now()
	for idx, it := range items {
		rk := Rank{
			SeerID:       it.SeerID,
			RatingSum:    it.RatingSum,
			ReviewCount:  it.ReviewCount,
			RankPosition: idx + 1,
			UpdatedAt:    now,
		}
		_, _, _ = s.rankCol.Add(ctx, rk)
	}
	return nil
}

func (s *RankService) CalculateCommission(ctx context.Context, month string, percent float64) error {
	q := s.paymentCol.Where("status", "==", "paid")
	docs, err := q.Documents(ctx).GetAll()
	if err != nil {
		return err
	}
	type payAgg struct {
		total float64
	}
	aggMap := make(map[string]*payAgg)
	for _, doc := range docs {
		var p Payment
		if err := doc.DataTo(&p); err != nil {
			continue
		}
		if p.CreatedAt.Format("2006-01") != month {
			continue
		}
		seerID := p.ProviderRefID
		if _, ok := aggMap[seerID]; !ok {
			aggMap[seerID] = &payAgg{}
		}
		aggMap[seerID].total += float64(p.Amount)
	}
	oldDocs, _ := s.commissionCol.Where("month", "==", month).Documents(ctx).GetAll()
	for _, od := range oldDocs {
		od.Ref.Delete(ctx)
	}
	now := time.Now()
	for seerID, a := range aggMap {
		c := Commission{
			SeerID:    seerID,
			Amount:    a.total * percent / 100.0,
			Month:     month,
			CreatedAt: now,
		}
		_, _, _ = s.commissionCol.Add(ctx, c)
	}
	return nil
}

func (s *RankService) CalculateBonus(ctx context.Context, period string, topN int, bonusAmount float64) error {
	cutoff := time.Now().AddDate(0, 0, -90)
	q := s.rankCol.Where("updatedAt", ">=", cutoff)
	docs, err := q.Documents(ctx).GetAll()
	if err != nil {
		return err
	}
	type rItem struct {
		SeerID       string
		RankPosition int
	}
	var items []rItem
	for _, doc := range docs {
		var r Rank
		if err := doc.DataTo(&r); err != nil {
			continue
		}
		items = append(items, rItem{SeerID: r.SeerID, RankPosition: r.RankPosition})
	}
	// Sort by RankPosition asc
	for i := 0; i < len(items); i++ {
		for j := i + 1; j < len(items); j++ {
			if items[j].RankPosition < items[i].RankPosition {
				items[i], items[j] = items[j], items[i]
			}
		}
	}
	oldDocs, _ := s.bonusCol.Where("period", "==", period).Documents(ctx).GetAll()
	for _, od := range oldDocs {
		od.Ref.Delete(ctx)
	}
	now := time.Now()
	count := 0
	for _, it := range items {
		if count >= topN {
			break
		}
		b := Bonus{
			SeerID:    it.SeerID,
			Amount:    bonusAmount,
			Period:    period,
			CreatedAt: now,
		}
		_, _, _ = s.bonusCol.Add(ctx, b)
		count++
	}
	return nil
}
