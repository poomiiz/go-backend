// go-backend/internal/routes/ai.go
package routes

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/poomiiz/go-backend/internal/services"
	"github.com/poomiiz/go-backend/internal/utils"
)

// RegisterAIRoutes ลงทะเบียน /ai/xxxx
func RegisterAIRoutes(r *gin.Engine) {
	aiClient := services.NewAIServiceClient(os.Getenv("AI_ROUTER_URL"))

	// /ai/interpret
	r.POST("/ai/interpret", func(c *gin.Context) {
		var req services.AIInterpretRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		aiResp, err := aiClient.Interpret(ctx, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		utils.SaveInterpretResult(req.UserID, req.ConversationID, aiResp.Intent, aiResp.Confidence)
		c.JSON(http.StatusOK, aiResp)
	})

	// /ai/summarize
	r.POST("/ai/summarize", func(c *gin.Context) {
		var req services.AISummarizeRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		aiResp, err := aiClient.Summarize(ctx, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, aiResp)
	})

	// /ai/chat
	r.POST("/ai/chat", func(c *gin.Context) {
		var req services.AIChatRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		utils.SaveUserMessage(req.UserID, req.ConversationID, req.Message)

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		aiResp, err := aiClient.Chat(ctx, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		utils.SaveBotMessage(req.UserID, req.ConversationID, aiResp.Response, aiResp.ModelUsed)
		c.JSON(http.StatusOK, aiResp)
	})

	// /ai/tune_prompt
	r.POST("/ai/tune_prompt", func(c *gin.Context) {
		var body struct {
			TuneID          string `json:"tuneId"`
			Model           string `json:"model"`
			CandidatePrompt string `json:"candidatePrompt"`
			TestQuestion    string `json:"testQuestion"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		result, err := aiClient.TunePrompt(ctx, body.TuneID, body.Model, body.CandidatePrompt, body.TestQuestion)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		utils.SavePromptTuneResult(body.TuneID, body.Model, body.CandidatePrompt, result)
		c.JSON(http.StatusOK, result)
	})
}
