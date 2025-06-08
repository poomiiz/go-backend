package routes

import (
	"context"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
)

type PromptConfig struct {
	Model       string    `json:"model"`
	Prompt      string    `json:"prompt"`
	Temperature float64   `json:"temperature"`
	MaxTokens   int       `json:"max_tokens"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func RegisterPromptRoutes(r *gin.Engine, firestoreClient *firestore.Client) {
	r.GET("/admin/prompts", func(c *gin.Context) {
		var prompts []map[string]interface{}
		iter := firestoreClient.Collection("configs").Documents(context.Background())
		for {
			doc, err := iter.Next()
			if err != nil {
				break
			}
			data := doc.Data()
			data["id"] = doc.Ref.ID
			prompts = append(prompts, data)
		}
		c.JSON(http.StatusOK, prompts)
	})

	r.POST("/admin/prompts/:key", func(c *gin.Context) {
		var body PromptConfig
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		body.UpdatedAt = time.Now()
		_, err := firestoreClient.Collection("configs").Doc(c.Param("key")).Set(context.Background(), body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "updated"})
	})
}
