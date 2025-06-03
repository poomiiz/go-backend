// internal/routes/booking.go
package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterBookingRoutes(r *gin.Engine) {
	grp := r.Group("/booking")
	{
		// ตัวอย่าง stub สำหรับสร้างจองคิว (ยังไม่ได้ implement จริง)
		grp.POST("/create", func(c *gin.Context) {
			c.JSON(http.StatusNotImplemented, gin.H{"error": "booking:create not implemented"})
		})
		// ตัวอย่าง stub สำหรับเลือก slot
		grp.POST("/select_slot", func(c *gin.Context) {
			c.JSON(http.StatusNotImplemented, gin.H{"error": "booking:select_slot not implemented"})
		})
		// ตัวอย่าง stub สำหรับ broadcast แจ้งเตือน (หากใช้)
		grp.GET("/notify", func(c *gin.Context) {
			c.JSON(http.StatusNotImplemented, gin.H{"error": "booking:notify not implemented"})
		})
	}
}
