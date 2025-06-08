// internal/routes/line_webhook.go
package routes

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/poomiiz/go-backend/internal/services"
	"github.com/poomiiz/go-backend/internal/utils"
)

// struct ที่ใช้รับ event จาก LINE webhook
type lineEvent struct {
	Events []eventObj `json:"events"`
}

type eventObj struct {
	ReplyToken string `json:"replyToken"`
	Message    struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"message"`
	Source struct {
		UserID string `json:"userId"`
	} `json:"source"`
}

func RegisterLineWebhook(r *gin.Engine) {
	aiURL := os.Getenv("AI_ROUTER_URL")
	aiModel := os.Getenv("AI_DEFAULT_MODEL")
	if aiModel == "" {
		aiModel = "gpt-4o"
	}

	r.POST("/webhook", func(c *gin.Context) {
		bodyBytes, _ := ioutil.ReadAll(c.Request.Body)
		var event lineEvent
		if err := json.Unmarshal(bodyBytes, &event); err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		for _, e := range event.Events {
			if e.Message.Type != "text" {
				continue
			}
			userID := e.Source.UserID
			replyToken := e.ReplyToken
			incomingText := e.Message.Text

			// สร้าง session ใหม่
			sessionId := uuid.New().String()

			// บันทึกข้อความผู้ใช้
			utils.SaveUserMessage(sessionId, userID, incomingText)

			// เรียก AI service
			aiReq := services.AIChatRequest{
				UserID:         userID,
				ConversationID: sessionId,
				Message:        incomingText,
				Model:          aiModel,
			}
			reqBody, _ := json.Marshal(aiReq)
			aiEndpoint := fmt.Sprintf("%s/chat", aiURL)
			httpReq, _ := http.NewRequestWithContext(context.Background(), "POST", aiEndpoint, bytes.NewBuffer(reqBody))
			httpReq.Header.Set("Content-Type", "application/json")
			httpClient := &http.Client{Timeout: 20 * time.Second}
			resp, err := httpClient.Do(httpReq)
			if err != nil {
				replyMessage(replyToken, "ขออภัย เกิดข้อผิดพลาด")
				continue
			}
			respBody, _ := ioutil.ReadAll(resp.Body)
			resp.Body.Close()

			if resp.StatusCode >= 400 {
				replyMessage(replyToken, "ขออภัย AI ไม่ตอบกลับ")
				continue
			}
			var aiResp services.AIChatResponse
			if err := json.Unmarshal(respBody, &aiResp); err != nil {
				replyMessage(replyToken, "เกิดข้อผิดพลาดในการประมวลผล")
				continue
			}

			// บันทึกข้อความจาก bot
			utils.SaveBotMessage(sessionId, userID, aiResp.Response, aiResp.ModelUsed)

			// ส่งข้อความกลับ LINE
			replyMessage(replyToken, aiResp.Response)

			// 🔁 สรุปบทสนทนา async
			go summarizeSession(sessionId)
		}

		c.Status(http.StatusOK)
	})
}

// ฟังก์ชันส่งข้อความตอบกลับ LINE
func replyMessage(replyToken string, message string) {
	endpoint := "https://api.line.me/v2/bot/message/reply"
	payload := map[string]interface{}{
		"replyToken": replyToken,
		"messages": []map[string]string{
			{
				"type": "text",
				"text": message,
			},
		},
	}
	jsonBody, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("LINE_CHANNEL_ACCESS_TOKEN"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("LINE API error:", err)
		return
	}
	defer resp.Body.Close()
}

// summarizeSession ดึงข้อความทั้งหมดจาก session แล้วส่งไปสรุป
func summarizeSession(sessionId string) {
	messages, _ := utils.GetSessionMessages(sessionId)
	fullText := utils.JoinText(messages)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	summary, err1 := services.AISummarize(ctx, fullText)
	intent, emotion, err2 := services.AIInterpret(ctx, fullText)

	if err1 == nil && err2 == nil {
		utils.SaveSummary(sessionId, summary, intent, fmt.Sprintf("%.2f", emotion))
	}

}
