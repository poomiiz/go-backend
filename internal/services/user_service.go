package services

import (
	"context"
	"errors"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/poomiiz/go-backend/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Email        string    `firestore:"email"`
	PasswordHash string    `firestore:"passwordHash"`
	Role         string    `firestore:"role"`
	TwoFASecret  string    `firestore:"twoFASecret"`
	Enabled2FA   bool      `firestore:"enabled2FA"`
	CreatedAt    time.Time `firestore:"createdAt"`
	UpdatedAt    time.Time `firestore:"updatedAt"`
}

type UserService struct {
	col *firestore.CollectionRef
}

func NewUserService() *UserService {
	return &UserService{
		col: utils.Client.Collection("users"),
	}
}

func (s *UserService) Register(ctx context.Context, email, plainPassword, role string) (string, error) {
	// ตรวจซ้ำ email
	q := s.col.Where("email", "==", email).Limit(1)
	docs, err := q.Documents(ctx).GetAll()
	if err != nil {
		return "", err
	}
	if len(docs) > 0 {
		return "", errors.New("email already registered")
	}

	// Hash password
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	now := time.Now()
	u := User{
		Email:        email,
		PasswordHash: string(hashBytes),
		Role:         role,
		TwoFASecret:  "",
		Enabled2FA:   false,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	docRef, _, err := s.col.Add(ctx, u)
	if err != nil {
		return "", err
	}
	return docRef.ID, nil
}

func (s *UserService) Login(ctx context.Context, email, plainPassword string) (string, *User, error) {
	q := s.col.Where("email", "==", email).Limit(1)
	docs, err := q.Documents(ctx).GetAll()
	if err != nil {
		return "", nil, err
	}
	if len(docs) == 0 {
		return "", nil, errors.New("invalid credentials")
	}
	var u User
	if err := docs[0].DataTo(&u); err != nil {
		return "", nil, err
	}
	if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(plainPassword)) != nil {
		return "", nil, errors.New("invalid credentials")
	}
	return docs[0].Ref.ID, &u, nil
}

func (s *UserService) GetByID(ctx context.Context, userID string) (*User, error) {
	doc, err := s.col.Doc(userID).Get(ctx)
	if err != nil {
		return nil, err
	}
	var u User
	if err := doc.DataTo(&u); err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *UserService) EnableTwoFA(ctx context.Context, userID, secret string) error {
	_, err := s.col.Doc(userID).Update(ctx, []firestore.Update{
		{Path: "twoFASecret", Value: secret},
		{Path: "enabled2FA", Value: true},
		{Path: "updatedAt", Value: time.Now()},
	})
	return err
}

func (s *UserService) DisableTwoFA(ctx context.Context, userID string) error {
	_, err := s.col.Doc(userID).Update(ctx, []firestore.Update{
		{Path: "twoFASecret", Value: ""},
		{Path: "enabled2FA", Value: false},
		{Path: "updatedAt", Value: time.Now()},
	})
	return err
}

func (s *UserService) UpdateProfile(ctx context.Context, userID, newEmail, newRole string) error {
	_, err := s.col.Doc(userID).Update(ctx, []firestore.Update{
		{Path: "email", Value: newEmail},
		{Path: "role", Value: newRole},
		{Path: "updatedAt", Value: time.Now()},
	})
	return err
}
