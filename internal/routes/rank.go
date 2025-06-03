package routes

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/poomiiz/go-backend/internal/services"
)

func RegisterRankRoutes(r *gin.Engine) {
	rankSvc := services.NewRankService()

	grp := r.Group("/rank")
	{
		grp.POST("/calc", func(c *gin.Context) {
			var payload struct {
				Since   string  `json:"since"` // timestamp ในรูป "2025-06-01T00:00:00Z"
				Percent float64 `json:"percent"`
			}
			if err := c.ShouldBindJSON(&payload); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
				return
			}
			sinceTime, err := time.Parse(time.RFC3339, payload.Since)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format"})
				return
			}
			err = rankSvc.CalculateRankings(c.Request.Context(), sinceTime)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			err = rankSvc.CalculateCommission(c.Request.Context(), sinceTime.Format("2006-01"), payload.Percent)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			period := time.Now().Format("2006-Q1") // ปรับตาม logic quarter จริง
			err = rankSvc.CalculateBonus(c.Request.Context(), period, 3, 1000.0)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"status": "ranking calculated"})
		})
	}
}
