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
	"github.com/google/uuid" // ใช้สำหรับสร้าง sessionId แบบสุ่ม
	"github.com/poomiiz/go-backend/internal/services"
	"github.com/poomiiz/go-backend/internal/utils"
)

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

			// **1) สร้าง session ใหม่ทุกครั้ง (หรืออาจจะเช็คว่า user มี session ค้างไว้หรือไม่)
			sessionId := uuid.New().String() // ex: "e4b8a3cd-9f7b-4d9a-8f1d-3c1234567890"

			// 2) บันทึกข้อความจาก user ลง Firestore ใช้ sessionId
			utils.SaveUserMessage(sessionId, userID, incomingText)

			// 3) เรียก AI Service แล้วรับผลลัพธ์
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

			// 4) บันทึกข้อความ bot ลง Firestore ใช้ sessionId เดิม
			utils.SaveBotMessage(sessionId, userID, aiResp.Response, aiResp.ModelUsed)

			// 5) ส่ง reply กลับ LINE
			replyMessage(replyToken, aiResp.Response)
		}

		c.Status(http.StatusOK)
	})
}
