package utils

import (
	"context"
	"log"
	"time"
)

// SaveUserMessage บันทึกข้อความจาก user ลง "conversations/{convId}/messages"
func SaveUserMessage(userID, convID, text string) {
	ctx := context.Background()
	doc := Client.Collection("conversations").Doc(convID).Collection("messages").NewDoc()
	payload := map[string]interface{}{
		"sender":    "user",
		"text":      text,
		"modelUsed": "",
		"timestamp": time.Now(),
	}
	_, err := doc.Set(ctx, payload)
	if err != nil {
		log.Println("Error saving user message:", err)
	}
}

// SaveBotMessage บันทึกข้อความจาก bot
func SaveBotMessage(userID, convID, text, modelUsed string) {
	ctx := context.Background()
	doc := Client.Collection("conversations").Doc(convID).Collection("messages").NewDoc()
	payload := map[string]interface{}{
		"sender":    "bot",
		"text":      text,
		"modelUsed": modelUsed,
		"timestamp": time.Now(),
	}
	_, err := doc.Set(ctx, payload)
	if err != nil {
		log.Println("Error saving bot message:", err)
	}
}

// SaveInterpretResult บันทึก intent/emotion analysis
func SaveInterpretResult(userID, convID, intent string, confidence float64) {
	ctx := context.Background()
	doc := Client.Collection("conversations").Doc(convID).Collection("interpretations").NewDoc()
	payload := map[string]interface{}{
		"intent":     intent,
		"confidence": confidence,
		"timestamp":  time.Now(),
	}
	_, err := doc.Set(ctx, payload)
	if err != nil {
		log.Println("Error saving interpret result:", err)
	}
}

// SavePromptTuneResult บันทึกผล Prompt Tuning
func SavePromptTuneResult(tuneID, model, prompt string, result map[string]interface{}) {
	ctx := context.Background()
	doc := Client.Collection("prompt_tunes").Doc(tuneID).Collection("variants").Doc(model)
	payload := map[string]interface{}{
		"model":           model,
		"promptText":      prompt,
		"generatedResult": result,
		"lastTestedAt":    time.Now(),
	}
	_, err := doc.Set(ctx, payload)
	if err != nil {
		log.Println("Error saving prompt tune result:", err)
	}
}
