// go-backend/internal/services/ai_router_service.go
package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// (อย่า import "github.com/poomiiz/go-backend/internal/routes" ที่นี่เด็ดขาด)

type AIInterpretRequest struct {
	UserID         string `json:"userId"`
	ConversationID string `json:"conversationId"`
	Message        string `json:"message"`
}
type AIInterpretResponse struct {
	Intent     string  `json:"intent"`
	Confidence float64 `json:"confidence"`
}

type AISummarizeRequest struct {
	ConversationID string   `json:"conversationId"`
	Messages       []string `json:"messages"`
}
type AISummarizeResponse struct {
	Summary string `json:"summary"`
}

type AIChatRequest struct {
	UserID         string `json:"userId"`
	ConversationID string `json:"conversationId"`
	Message        string `json:"message"`
	Model          string `json:"model"`
}
type AIChatResponse struct {
	Response        string  `json:"response"`
	ModelUsed       string  `json:"modelUsed"`
	ConfidenceScore float64 `json:"confidenceScore"`
	Summary         string  `json:"summary"`
}

type AIServiceClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewAIServiceClient(baseURL string) *AIServiceClient {
	return &AIServiceClient{
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: 15 * time.Second},
	}
}

// Interpret
func (c *AIServiceClient) Interpret(ctx context.Context, req AIInterpretRequest) (*AIInterpretResponse, error) {
	url := fmt.Sprintf("%s/interpret", c.baseURL)
	data, _ := json.Marshal(req)
	httpReq, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(data))
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("ai-service interpret status %d", resp.StatusCode)
	}
	var aiResp AIInterpretResponse
	if err := json.NewDecoder(resp.Body).Decode(&aiResp); err != nil {
		return nil, err
	}
	return &aiResp, nil
}

// Summarize
func (c *AIServiceClient) Summarize(ctx context.Context, req AISummarizeRequest) (*AISummarizeResponse, error) {
	url := fmt.Sprintf("%s/summarize", c.baseURL)
	data, _ := json.Marshal(req)
	httpReq, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(data))
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("ai-service summarize status %d", resp.StatusCode)
	}
	var aiResp AISummarizeResponse
	if err := json.NewDecoder(resp.Body).Decode(&aiResp); err != nil {
		return nil, err
	}
	return &aiResp, nil
}

// Chat
func (c *AIServiceClient) Chat(ctx context.Context, req AIChatRequest) (*AIChatResponse, error) {
	url := fmt.Sprintf("%s/chat", c.baseURL)
	data, _ := json.Marshal(req)
	httpReq, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(data))
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("ai-service chat status %d", resp.StatusCode)
	}
	var aiResp AIChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&aiResp); err != nil {
		return nil, err
	}
	return &aiResp, nil
}

// TunePrompt
func (c *AIServiceClient) TunePrompt(ctx context.Context, tuneID, model, candidatePrompt, testQuestion string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/tune_prompt", c.baseURL)
	payload := map[string]string{
		"tuneId":          tuneID,
		"model":           model,
		"candidatePrompt": candidatePrompt,
		"testQuestion":    testQuestion,
	}
	data, _ := json.Marshal(payload)
	httpReq, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(data))
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("ai-service tune_prompt status %d", resp.StatusCode)
	}
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}
