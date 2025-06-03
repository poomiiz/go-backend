package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/poomiiz/go-backend/internal/services"
)

func RegisterReviewRoutes(r *gin.Engine) {
	reviewSvc := services.NewReviewService()
	grp := r.Group("/review")
	{
		grp.POST("/submit", func(c *gin.Context) {
			var payload struct {
				UserID  string `json:"userId"`
				SeerID  string `json:"seerId"`
				Rating  int    `json:"rating"`
				Content string `json:"content"`
			}
			if err := c.ShouldBindJSON(&payload); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
				return
			}
			revID, err := reviewSvc.SubmitReview(c.Request.Context(), payload.UserID, payload.SeerID, payload.Rating, payload.Content)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusCreated, gin.H{"reviewId": revID})
		})

		grp.GET("/pending", func(c *gin.Context) {
			list, err := reviewSvc.GetPendingReviews(c.Request.Context())
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, list)
		})

		grp.POST("/approve", func(c *gin.Context) {
			var payload struct {
				ReviewID string `json:"reviewId"`
			}
			if err := c.ShouldBindJSON(&payload); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
				return
			}
			err := reviewSvc.ApproveReview(c.Request.Context(), payload.ReviewID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"status": "approved"})
		})

		grp.POST("/reject", func(c *gin.Context) {
			var payload struct {
				ReviewID string `json:"reviewId"`
			}
			if err := c.ShouldBindJSON(&payload); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
				return
			}
			err := reviewSvc.RejectReview(c.Request.Context(), payload.ReviewID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"status": "rejected"})
		})

		grp.POST("/appeal", func(c *gin.Context) {
			var payload struct {
				ReviewID string `json:"reviewId"`
				UserID   string `json:"userId"`
				Reason   string `json:"reason"`
			}
			if err := c.ShouldBindJSON(&payload); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
				return
			}
			appID, err := reviewSvc.AppealReview(c.Request.Context(), payload.ReviewID, payload.UserID, payload.Reason)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusCreated, gin.H{"appealId": appID})
		})
	}
}
