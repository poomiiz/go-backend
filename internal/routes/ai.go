// internal/routes/ai.go
package routes

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/poomiiz/go-backend/internal/services"
)

// RegisterAIRoutes ลงทะเบียน endpoint กลุ่ม /ai
func RegisterAIRoutes(r *gin.Engine) {
	aiSvc := services.NewAIRouterService(os.Getenv("AI_ROUTER_URL")) // ต้องตั้ง AI_ROUTER_URL ใน .env

	grp := r.Group("/ai")
	{
		// POST /ai/chat
		grp.POST("/chat", func(c *gin.Context) {
			var payload struct {
				UserID string `json:"user_id"`
				Prompt string `json:"prompt"`
				Model  string `json:"model"`
			}
			if err := c.ShouldBindJSON(&payload); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
				return
			}
			ctx := c.Request.Context()
			resp, err := aiSvc.Chat(ctx, services.AIRequest{
				UserID: payload.UserID,
				Prompt: payload.Prompt,
				Model:  payload.Model,
			})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, resp)
		})

		// POST /ai/daily_card
		grp.POST("/daily_card", func(c *gin.Context) {
			var payload struct {
				UserID string `json:"user_id"`
				Deck   string `json:"deck"`
				Date   string `json:"date"`
			}
			if err := c.ShouldBindJSON(&payload); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
				return
			}
			ctx := c.Request.Context()
			resp, err := aiSvc.DailyCard(ctx, payload.UserID, payload.Deck, payload.Date)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, resp)
		})

		// POST /ai/interpret
		grp.POST("/interpret", func(c *gin.Context) {
			var payload struct {
				UserID  string   `json:"user_id"`
				CardIDs []string `json:"cards"`
			}
			if err := c.ShouldBindJSON(&payload); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
				return
			}
			ctx := c.Request.Context()
			resp, err := aiSvc.InterpretCard(ctx, payload.UserID, payload.CardIDs)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, resp)
		})
	}
}
