// internal/utils/firestore_client.go
package utils

import (
	"context"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/firestore"
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
