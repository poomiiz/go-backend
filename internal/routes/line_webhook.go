package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// โครงสร้างเบื้องต้นของ JSON payload ที่ LINE ส่งมา
type lineEvent struct {
	Events []struct {
		ReplyToken string `json:"replyToken"`
		Source     struct {
			UserID string `json:"userId"`
			Type   string `json:"type"`
		} `json:"source"`
		Message struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"message"`
	} `json:"events"`
}

// ฟังก์ชันส่งข้อความตอบกลับ (Reply)
func replyMessage(replyToken, text string) error {
	channelToken := os.Getenv("LINE_CHANNEL_TOKEN")
	if channelToken == "" {
		return fmt.Errorf("LINE_CHANNEL_TOKEN is not set")
	}

	payload := map[string]interface{}{
		"replyToken": replyToken,
		"messages": []map[string]string{
			{
				"type": "text",
				"text": text,
			},
		},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	url := "https://api.line.me/v2/bot/message/reply"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+channelToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		respBody, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("LINE reply API error: status %d, body %s", resp.StatusCode, string(respBody))
	}
	return nil
}

// RegisterLineWebhook ลงทะเบียน endpoint /webhook ให้พิมพ์ payload และตอบกลับอัตโนมัติ
func RegisterLineWebhook(r *gin.Engine) {
	r.POST("/webhook", func(c *gin.Context) {
		bodyBytes, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			fmt.Println("Error reading request body:", err)
			c.Status(http.StatusInternalServerError)
			return
		}
		fmt.Println("LINE webhook payload:", string(bodyBytes))

		var event lineEvent
		if err := json.Unmarshal(bodyBytes, &event); err != nil {
			fmt.Println("Error parsing JSON:", err)
			c.Status(http.StatusBadRequest)
			return
		}

		for _, e := range event.Events {
			if e.Message.Type == "text" {
				incoming := e.Message.Text
				replyToken := e.ReplyToken
				// ลบการประกาศ userID ที่ไม่ได้ใช้งานออก
				// userID := e.Source.UserID

				replyText := fmt.Sprintf("คุณพิมพ์ว่า: %s", incoming)

				if err := replyMessage(replyToken, replyText); err != nil {
					fmt.Println("Error replying message:", err)
				}
			}
		}

		c.Status(http.StatusOK)
	})
}
