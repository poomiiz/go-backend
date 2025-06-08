// internal/utils/firestore_client.go
package utils

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

// Client ถูก initialize ไว้ใน InitFirestore()
var Client *firestore.Client

// ----------------------------------------------------------------------------
// SaveUserMessage – บันทึกข้อความของผู้ใช้ ลง subcollection "messages_YYYY_MM"
func SaveUserMessage(sessionId, userId, text string) {
	ctx := context.Background()

	// 1. ตั้งชื่อ document ชั้นบนตาม sessionId
	docRef := Client.Collection("conversations").Doc(sessionId)

	// 2. สร้างหรืออัปเดต field "userId" ใน session document (ใช้ MergeAll เพื่อไม่ลบ field อื่น)
	_, err := docRef.Set(ctx, map[string]interface{}{
		"userId": userId,
	}, firestore.MergeAll)
	if err != nil {
		log.Println("Error setting conversation doc:", err)
		return
	}

	// 3. หาเดือน/ปีปัจจุบัน เพื่อจัด subcollection partition
	now := time.Now()
	year, month, _ := now.Date()
	subcol := fmt.Sprintf("messages_%04d_%02d", year, int(month))
	//    → example: "messages_2025_06"

	// 4. สร้าง document ใน subcollection partition นั้น
	msgRef := docRef.Collection(subcol).NewDoc()
	payload := map[string]interface{}{
		"sender":    "user",
		"text":      text,
		"modelUsed": "",
		"timestamp": now,
		// ใส่ expireAt สำหรับ TTL (ถ้าใช้ TTL)
		"expireAt": now.Add(24 * time.Hour),
	}
	if _, err := msgRef.Set(ctx, payload); err != nil {
		log.Println("Error saving user message:", err)
		return
	}
}

// ----------------------------------------------------------------------------
// SaveBotMessage – บันทึกข้อความของ bot ลง subcollection "messages_YYYY_MM"
func SaveBotMessage(sessionId, userId, text, modelUsed string) {
	ctx := context.Background()

	// Document ชั้นบน
	docRef := Client.Collection("conversations").Doc(sessionId)

	// (ไม่จำเป็นต้องตั้ง userId ซ้ำ เพราะ SaveUserMessage หมายถึง session จะถูกสร้างครั้งแรกที่ user ส่งข้อความ)
	// แต่หากต้องการก็ทำ MergeAll ได้ไม่เสียหาย:
	_, _ = docRef.Set(ctx, map[string]interface{}{
		"userId": userId,
	}, firestore.MergeAll)

	now := time.Now()
	year, month, _ := now.Date()
	subcol := fmt.Sprintf("messages_%04d_%02d", year, int(month)) // "messages_2025_06"

	msgRef := docRef.Collection(subcol).NewDoc()
	payload := map[string]interface{}{
		"sender":    "bot",
		"text":      text,
		"modelUsed": modelUsed,
		"timestamp": now,
		"expireAt":  now.Add(24 * time.Hour),
	}
	if _, err := msgRef.Set(ctx, payload); err != nil {
		log.Println("Error saving bot message:", err)
		return
	}
}

// GetSessionMessages ดึงข้อความ user ทั้งหมดใน session
func GetSessionMessages(sessionID string) ([]string, error) {
	ctx := context.Background()
	docRef := Client.Collection("conversations").Doc(sessionID)

	now := time.Now()
	year, month, _ := now.Date()
	subcol := fmt.Sprintf("messages_%04d_%02d", year, int(month))

	iter := docRef.Collection(subcol).Where("sender", "==", "user").Documents(ctx)
	defer iter.Stop()

	var messages []string
	for {
		doc, err := iter.Next()
		if err != nil {
			break
		}
		text, _ := doc.Data()["text"].(string)
		messages = append(messages, text)
	}
	return messages, nil
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
func InitFirestore() error {
	projectID := os.Getenv("FIREBASE_PROJECT_ID")
	cred := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")

	log.Println("📌 FIREBASE_PROJECT_ID =", projectID)
	log.Println("📌 GOOGLE_APPLICATION_CREDENTIALS =", cred)

	if projectID == "" {
		return fmt.Errorf("FIREBASE_PROJECT_ID is not set")
	}
	if cred == "" {
		return fmt.Errorf("GOOGLE_APPLICATION_CREDENTIALS is not set")
	}

	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("firestore.NewClient error: %w", err)
	}

	Client = client
	log.Println("✅ Firestore client initialized")
	return nil
}
func SaveSummary(sessionId, summary, intent, emotion string) {
	ctx := context.Background()
	_, err := Client.Collection("conversations").Doc(sessionId).Set(ctx, map[string]interface{}{
		"summary":   summary,
		"intent":    intent,
		"emotion":   emotion,
		"summaryAt": time.Now(),
	}, firestore.MergeAll)
	if err != nil {
		log.Println("Error saving summary:", err)
	}
}

type Conversation struct {
	ID   string
	Data map[string]interface{}
}

func JoinText(lines []string) string {
	return strings.Join(lines, "\n")
}
func QueryUserConversations(userID string) []Conversation {
	ctx := context.Background()
	iter := Client.Collection("conversations").Where("userId", "==", userID).Documents(ctx)

	result := []Conversation{}
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			continue
		}
		result = append(result, Conversation{
			ID:   doc.Ref.ID,
			Data: doc.Data(),
		})
	}
	return result
}

func CloseFirestore() {
	if Client != nil {
		if err := Client.Close(); err != nil {
			log.Println("Error closing Firestore:", err)
		}
	}
}

// GetFirestoreClient คืนค่า Firestore client ที่ถูก init ไว้แล้ว
func GetFirestoreClient() *firestore.Client {
	return Client
}
