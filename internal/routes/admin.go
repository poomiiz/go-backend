// internal/routes/admin.go
package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ---------- Stub Middleware & Handlers ----------

// authMiddleware: ถ้าไม่ต้องตรวจสิทธิ์จริง ให้เขียนเป็น stub แบบนี้ก่อน
func authMiddleware(c *gin.Context) {
	// ถ้ามี logic ตรวจ token / session ให้ใส่ที่นี่
	// ตอนนี้แค่ผ่านไปเลย
	c.Next()
}

// getPromptTuneHandler: stub ดึง Prompt Template จาก Firestore
func getPromptTuneHandler(c *gin.Context) {
	tuneId := c.Param("tuneId")
	// TODO: เรียก utils.GetPromptTune(tuneId) หรือ service ดึงข้อมูลจริงมา
	// สำหรับ stub ให้ return JSON เปล่าๆ หรือโครงสร้างตามที่คาด
	c.JSON(http.StatusOK, gin.H{
		"tuneId":   tuneId,
		"variants": []interface{}{}, // ส่ง list เปล่าไปก่อน
	})
}

// approvePromptVariantHandler: stub บันทึกการอนุมัติ prompt variant
func approvePromptVariantHandler(c *gin.Context) {
	tuneId := c.Param("tuneId")
	var body struct {
		Model      string `json:"model"`
		PromptText string `json:"promptText"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// TODO: บันทึกลง Firestore จริง
	c.JSON(http.StatusOK, gin.H{
		"status": "approved",
		"tuneId": tuneId,
		"model":  body.Model,
	})
}

// listConversationsHandler: stub ดึงรายชื่อทั้งหมดของ conversations
func listConversationsHandler(c *gin.Context) {
	// TODO: ดึง list conversationId จาก Firestore
	c.JSON(http.StatusOK, gin.H{
		"conversations": []string{}, // ส่ง list เปล่าไปก่อน
	})
}

// getMessagesHandler: stub ดึงข้อความทั้งหมดใน conversation หนึ่ง
func getMessagesHandler(c *gin.Context) {
	convId := c.Param("convId")
	// TODO: ดึง /conversations/{convId}/messages จาก Firestore
	c.JSON(http.StatusOK, gin.H{
		"conversationId": convId,
		"messages":       []interface{}{}, // ส่ง list เปล่าไปก่อน
	})
}

// getInterpretationsHandler: stub ดึงผล intent/emotion ของ conversation
func getInterpretationsHandler(c *gin.Context) {
	convId := c.Param("convId")
	// TODO: ดึง /conversations/{convId}/interpretations จาก Firestore
	c.JSON(http.StatusOK, gin.H{
		"conversationId":  convId,
		"interpretations": []interface{}{}, // ส่ง list เปล่าไปก่อน
	})
}

// regenerateSummaryHandler: stub ให้ระบบสรุปบทสนทนาใหม่
func regenerateSummaryHandler(c *gin.Context) {
	convId := c.Param("convId")
	// TODO: เรียก AI Service /ai/summarize เพื่อสรุปบทสนทนาใหม่
	c.JSON(http.StatusOK, gin.H{
		"conversationId": convId,
		"summary":        "สรุปใหม่ (stub)",
	})
}

// ---------- RegisterAdminRoutes ----------

func RegisterAdminRoutes(r *gin.Engine) {
	// ใช้ authMiddleware (แม้เป็น stub) เพื่อป้องกัน route นี้ไม่ให้คนทั่วไปเข้าถึงได้
	admin := r.Group("/admin", authMiddleware)
	{
		admin.GET("/prompt_tunes/:tuneId", getPromptTuneHandler)
		admin.POST("/prompt_tunes/:tuneId/approve", approvePromptVariantHandler)
		admin.GET("/conversations", listConversationsHandler)
		admin.GET("/conversations/:convId/messages", getMessagesHandler)
		admin.GET("/conversations/:convId/interpretations", getInterpretationsHandler)
		admin.POST("/conversations/:convId/regenerate_summary", regenerateSummaryHandler)
	}
}
