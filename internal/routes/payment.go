package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/poomiiz/go-backend/internal/services"
)

func RegisterPaymentRoutes(r *gin.Engine) {
	paySvc := services.NewPaymentService(5.0) // หรือใส่ percent เป็น env

	grp := r.Group("/payment")
	{
		grp.POST("/create", func(c *gin.Context) {
			var payload struct {
				UserID        string `json:"userId"`
				Amount        int64  `json:"amount"`
				Provider      string `json:"provider"`
				ProviderRefID string `json:"providerRefId"`
			}
			if err := c.ShouldBindJSON(&payload); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
				return
			}
			payID, err := paySvc.CreatePayment(c.Request.Context(), payload.UserID, payload.Amount, payload.Provider, payload.ProviderRefID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusCreated, gin.H{"paymentId": payID})
		})

		grp.POST("/verify", func(c *gin.Context) {
			var payload struct {
				PaymentID string `json:"paymentId"`
			}
			if err := c.ShouldBindJSON(&payload); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
				return
			}
			err := paySvc.VerifyPayment(c.Request.Context(), payload.PaymentID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"status": "verified"})
		})
	}
}
