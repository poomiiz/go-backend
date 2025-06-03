package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// AIRequest โครงสร้างที่ส่งไปยัง ai-service
type AIRequest struct {
	UserID string                 `json:"user_id"`
	Prompt string                 `json:"prompt"`
	Model  string                 `json:"model"`
	Extra  map[string]interface{} `json:"extra,omitempty"`
}

// AIResponse โครงสร้างที่รับมาจาก ai-service
type AIResponse struct {
	Reply string `json:"reply"`
	// …สามารถขยายได้ตาม response ของ FastAPI
}

type AIRouterService struct {
	baseURL    string
	httpClient *http.Client
}

// NewAIRouterService สร้าง instance ด้วย baseURL เช่น "http://localhost:8000/ai"
func NewAIRouterService(baseURL string) *AIRouterService {
	return &AIRouterService{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Chat: ส่ง POST ไปที่ <baseURL>/chat
func (s *AIRouterService) Chat(ctx context.Context, reqPayload AIRequest) (*AIResponse, error) {
	url := fmt.Sprintf("%s/chat", s.baseURL)
	return s.sendRequest(ctx, url, reqPayload)
}

// DailyCard: ส่ง POST ไปที่ <baseURL>/daily_card
func (s *AIRouterService) DailyCard(ctx context.Context, userID, deck, date string) (*AIResponse, error) {
	payload := map[string]interface{}{
		"user_id": userID,
		"deck":    deck,
		"date":    date,
	}
	url := fmt.Sprintf("%s/daily_card", s.baseURL)
	return s.sendRaw(ctx, url, payload)
}

// InterpretCard: ส่ง POST ไปที่ <baseURL>/interpret
func (s *AIRouterService) InterpretCard(ctx context.Context, userID string, cardIDs []string) (*AIResponse, error) {
	payload := map[string]interface{}{
		"user_id": userID,
		"cards":   cardIDs,
	}
	url := fmt.Sprintf("%s/interpret", s.baseURL)
	return s.sendRaw(ctx, url, payload)
}

func (s *AIRouterService) sendRequest(ctx context.Context, url string, payload AIRequest) (*AIResponse, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("ai-service returned status %d", resp.StatusCode)
	}

	var aiResp AIResponse
	if err := json.NewDecoder(resp.Body).Decode(&aiResp); err != nil {
		return nil, err
	}
	return &aiResp, nil
}

func (s *AIRouterService) sendRaw(ctx context.Context, url string, payload map[string]interface{}) (*AIResponse, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("ai-service returned status %d", resp.StatusCode)
	}

	var aiResp AIResponse
	if err := json.NewDecoder(resp.Body).Decode(&aiResp); err != nil {
		return nil, err
	}
	return &aiResp, nil
}
