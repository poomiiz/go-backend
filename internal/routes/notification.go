package routes

import (
	"net/http"

	"os"

	"github.com/gin-gonic/gin"
	"github.com/poomiiz/go-backend/internal/services"
)

func RegisterNotificationRoutes(r *gin.Engine) {
	lineToken := os.Getenv("LINE_CHANNEL_TOKEN")
	telegramURL := os.Getenv("TELEGRAM_BOT_URL")
	telegramAuth := os.Getenv("TELEGRAM_BOT_AUTH") // ถ้ามี
	notifSvc := services.NewNotificationService(lineToken, telegramURL, telegramAuth)

	grp := r.Group("/notification")
	{
		grp.POST("/line", func(c *gin.Context) {
			var payload struct {
				To      string `json:"to"`
				Message string `json:"message"`
			}
			if err := c.ShouldBindJSON(&payload); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
				return
			}
			err := notifSvc.SendLineMessage(c.Request.Context(), payload.To, payload.Message)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		grp.POST("/telegram", func(c *gin.Context) {
			var payload struct {
				Type string                 `json:"type"`
				Data map[string]interface{} `json:"data"`
			}
			if err := c.ShouldBindJSON(&payload); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
				return
			}
			err := notifSvc.SendTelegramAlert(c.Request.Context(), payload.Type, payload.Data)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})
	}
}
