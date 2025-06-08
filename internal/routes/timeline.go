// internal/routes/timeline.go
package routes

import (
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/poomiiz/go-backend/internal/utils"
)

type TimelineItem struct {
	SessionID string `json:"sessionId"`
	Summary   string `json:"summary"`
	Intent    string `json:"intent"`
	Emotion   string `json:"emotion"`
	Timestamp string `json:"timestamp"`
}

func RegisterTimelineRoutes(r *gin.Engine) {
	r.GET("/timeline/:userId", func(c *gin.Context) {
		userID := c.Param("userId")
		sessions := utils.QueryUserConversations(userID)

		items := make([]TimelineItem, 0)
		for _, s := range sessions {
			items = append(items, TimelineItem{
				SessionID: s.ID,
				Summary:   s.Data["summary"].(string),
				Intent:    s.Data["intent"].(string),
				Emotion:   s.Data["emotion"].(string),
				Timestamp: s.Data["startedAt"].(string), // หรือใช้ time.Time.Format
			})
		}

		// เรียงตามเวลา (ใหม่ → เก่า)
		sort.Slice(items, func(i, j int) bool {
			return items[i].Timestamp > items[j].Timestamp
		})

		c.JSON(http.StatusOK, items)
	})
}
