// scripts/seed_ai_routing.go
package main

import (
	"context"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

func main() {
	ctx := context.Background()

	// โหลด service account
	sa := option.WithCredentialsFile("config/serviceAccountKey.json")
	client, err := firestore.NewClient(ctx, "moobiker-moomoon", sa)
	if err != nil {
		log.Fatalf("❌ Failed to connect Firestore: %v", err)
	}
	defer client.Close()

	// ✅ เพิ่ม temperature และ max_tokens สำหรับแต่ละ AI role
	data := map[string]map[string]interface{}{
		"chat": {
			"model":       "gpt-4o",
			"prompt_key":  "ai_prompt.chat",
			"temperature": 0.7,
			"max_tokens":  800,
			"updatedAt":   time.Now(),
		},
		"summarize": {
			"model":       "gpt-3.5-turbo",
			"prompt_key":  "ai_prompt.summarize",
			"temperature": 0.5,
			"max_tokens":  400,
			"updatedAt":   time.Now(),
		},
		"interpret": {
			"model":       "gpt-3.5-turbo",
			"prompt_key":  "ai_prompt.interpret",
			"temperature": 0.3,
			"max_tokens":  300,
			"updatedAt":   time.Now(),
		},
		"tarot": {
			"model":       "gpt-4o",
			"prompt_key":  "ai_prompt.tarot",
			"temperature": 0.8,
			"max_tokens":  1000,
			"updatedAt":   time.Now(),
		},
		"rag": {
			"model":       "meta-llama/Llama-3.3-70B-Instruct-Turbo-Free",
			"prompt_key":  "ai_prompt.rag",
			"temperature": 0.6,
			"max_tokens":  1200,
			"updatedAt":   time.Now(),
		},
		"test": {
			"model":       "gpt-4o",
			"prompt_key":  "ai_prompt.test",
			"temperature": 0.9,
			"max_tokens":  800,
			"updatedAt":   time.Now(),
		},
		"monitor": {
			"model":       "gpt-4o",
			"prompt_key":  "ai_prompt.monitor",
			"temperature": 0.0,
			"max_tokens":  500,
			"updatedAt":   time.Now(),
		},
	}

	for purpose, config := range data {
		_, err := client.Collection("ai_routing").Doc(purpose).Set(ctx, config, firestore.MergeAll) // ✅ MergeAll จะอัปเดตเฉพาะฟิลด์ ไม่ลบทิ้ง
		if err != nil {
			log.Printf("❌ Failed to write %s: %v", purpose, err)
		} else {
			log.Printf("✅ Seeded %s", purpose)
		}
	}
}
