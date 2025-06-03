package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type NotificationService struct {
	lineToken       string
	telegramBotURL  string // URL ของ telegram-alert-bot (เช่น "http://localhost:5000/alert")
	telegramBotAuth string // ถ้ามีใช้ Bearer token
}

func NewNotificationService(lineToken, telegramBotURL, telegramBotAuth string) *NotificationService {
	return &NotificationService{
		lineToken:       lineToken,
		telegramBotURL:  telegramBotURL,
		telegramBotAuth: telegramBotAuth,
	}
}

// SendLineMessage: ส่งข้อความไปยัง LINE Messaging API
func (s *NotificationService) SendLineMessage(ctx context.Context, toUserID, message string) error {
	type lineMessage struct {
		To       string `json:"to"`
		Messages []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"messages"`
	}

	payload := lineMessage{
		To: toUserID,
	}
	payload.Messages = []struct {
		Type string "json:\"type\""
		Text string "json:\"text\""
	}{
		{Type: "text", Text: message},
	}

	data, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.line.me/v2/bot/message/push", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.lineToken)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("line API returned status %d", resp.StatusCode)
	}
	return nil
}

// SendTelegramAlert: ส่ง alert ไปยัง telegram-alert-bot
func (s *NotificationService) SendTelegramAlert(ctx context.Context, alertType string, payload map[string]interface{}) error {
	body := map[string]interface{}{
		"type": alertType,
		"data": payload,
	}
	data, _ := json.Marshal(body)

	req, err := http.NewRequestWithContext(ctx, "POST", s.telegramBotURL, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if s.telegramBotAuth != "" {
		req.Header.Set("Authorization", "Bearer "+s.telegramBotAuth)
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("telegram-alert-bot returned status %d", resp.StatusCode)
	}
	return nil
}

// AlertCoinTopUp: ตัวช่วยส่ง alert กรณีเติมเหรียญ
func (s *NotificationService) AlertCoinTopUp(ctx context.Context, userID string, amount int64) error {
	payload := map[string]interface{}{
		"userId": userID,
		"amount": amount,
		"time":   time.Now().Format(time.RFC3339),
	}
	return s.SendTelegramAlert(ctx, "coin_topup", payload)
}

// AlertNewReview: ตัวช่วยส่ง alert กรณีมีรีวิวใหม่
func (s *NotificationService) AlertNewReview(ctx context.Context, userID string, reviewID string) error {
	payload := map[string]interface{}{
		"userId":   userID,
		"reviewId": reviewID,
		"time":     time.Now().Format(time.RFC3339),
	}
	return s.SendTelegramAlert(ctx, "new_review", payload)
}
