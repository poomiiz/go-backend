package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/poomiiz/go-backend/internal/services"
)

func RegisterCoinRoutes(r *gin.Engine) {
	coinSvc := services.NewCoinService()
	grp := r.Group("/coin")
	{
		grp.GET("/balance", func(c *gin.Context) {
			userID := c.Query("userId")
			bal, err := coinSvc.GetBalance(c.Request.Context(), userID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"balance": bal})
		})

		grp.POST("/topup", func(c *gin.Context) {
			var payload struct {
				UserID string `json:"userId"`
				Amount int64  `json:"amount"`
			}
			if err := c.ShouldBindJSON(&payload); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
				return
			}
			err := coinSvc.TopUp(c.Request.Context(), payload.UserID, payload.Amount)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		grp.POST("/transfer", func(c *gin.Context) {
			var payload struct {
				FromUserID string `json:"fromUserId"`
				ToUserID   string `json:"toUserId"`
				Amount     int64  `json:"amount"`
			}
			if err := c.ShouldBindJSON(&payload); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
				return
			}
			err := coinSvc.Transfer(c.Request.Context(), payload.FromUserID, payload.ToUserID, payload.Amount)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})
	}
}
