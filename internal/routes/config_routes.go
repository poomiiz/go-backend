// internal/routes/config_routes.go
package routes

import (
	"context"
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"github.com/poomiiz/go-backend/internal/utils" // 🔁 import ได้เพราะอยู่คนละ package
)

type PromptUpdateBody struct {
	Key    string `json:"key"`    // ex: ai_prompt.line
	Model  string `json:"model"`  // ex: gpt-4o
	Prompt string `json:"prompt"` // ex: actual prompt text
}

func RegisterAdminConfigRoutes(r *gin.Engine) {
	r.POST("/config/prompt/update", func(c *gin.Context) {
		var body struct {
			Key    string `json:"key"`
			Model  string `json:"model"`
			Prompt string `json:"prompt"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
			return
		}
		ctx := context.Background()
		_, err := utils.Client.Collection("configs").Doc(body.Key).Set(ctx, map[string]interface{}{
			"model":  body.Model,
			"prompt": body.Prompt,
		}, firestore.MergeAll)
		if err != nil {
			log.Println("❌ Failed to update prompt:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "update failed"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
}
