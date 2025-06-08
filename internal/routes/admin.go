// internal/routes/admin.go
package routes

import (
	"context"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"github.com/poomiiz/go-backend/internal/utils"
)

// ---------- Stub Middleware & Handlers ----------

func authMiddleware(c *gin.Context) {
	c.Next()
}

func getPromptTuneHandler(c *gin.Context) {
	tuneId := c.Param("tuneId")
	c.JSON(http.StatusOK, gin.H{
		"tuneId":   tuneId,
		"variants": []interface{}{},
	})
}

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
	c.JSON(http.StatusOK, gin.H{
		"status": "approved",
		"tuneId": tuneId,
		"model":  body.Model,
	})
}

func listConversationsHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"conversations": []string{},
	})
}

func getMessagesHandler(c *gin.Context) {
	convId := c.Param("convId")
	c.JSON(http.StatusOK, gin.H{
		"conversationId": convId,
		"messages":       []interface{}{},
	})
}

func getInterpretationsHandler(c *gin.Context) {
	convId := c.Param("convId")
	c.JSON(http.StatusOK, gin.H{
		"conversationId":  convId,
		"interpretations": []interface{}{},
	})
}

func regenerateSummaryHandler(c *gin.Context) {
	convId := c.Param("convId")
	c.JSON(http.StatusOK, gin.H{
		"conversationId": convId,
		"summary":        "สรุปใหม่ (stub)",
	})
}

// ---------- Register Admin Routes ----------

func RegisterAdminRoutes(r *gin.Engine) {
	admin := r.Group("/admin", authMiddleware)
	{
		admin.GET("/prompt_tunes/:tuneId", getPromptTuneHandler)
		admin.POST("/prompt_tunes/:tuneId/approve", approvePromptVariantHandler)
		admin.GET("/conversations", listConversationsHandler)
		admin.GET("/conversations/:convId/messages", getMessagesHandler)
		admin.GET("/conversations/:convId/interpretations", getInterpretationsHandler)
		admin.POST("/conversations/:convId/regenerate_summary", regenerateSummaryHandler)
		admin.POST("/config/prompt/update", func(c *gin.Context) {
			var payload struct {
				Key    string `json:"key"`
				Model  string `json:"model"`
				Prompt string `json:"prompt"`
			}
			if err := c.BindJSON(&payload); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
				return
			}
			if payload.Key == "" || payload.Model == "" || payload.Prompt == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "missing fields"})
				return
			}

			ctx := context.Background()
			doc := utils.Client.Collection("configs").Doc(payload.Key)
			_, err := doc.Set(ctx, map[string]interface{}{
				"model":     payload.Model,
				"prompt":    payload.Prompt,
				"updatedAt": time.Now(),
			}, firestore.MergeAll)

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "save failed"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"status": "updated"})
		})
	}
}
