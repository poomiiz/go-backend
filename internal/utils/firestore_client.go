package utils

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

var Client *firestore.Client

func InitFirestore() {
	ctx := context.Background()

	// อ่าน path ไปยัง Service Account JSON
	saPath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if saPath == "" {
		log.Fatal("GOOGLE_APPLICATION_CREDENTIALS is not set")
	}

	// อ่าน Project ID จาก env
	projectID := os.Getenv("FIREBASE_PROJECT_ID")
	if projectID == "" {
		log.Fatal("FIREBASE_PROJECT_ID is not set")
	}

	client, err := firestore.NewClient(ctx, projectID, option.WithCredentialsFile(saPath))
	if err != nil {
		log.Fatalf("Failed to create Firestore client: %v", err)
	}
	Client = client
}

func CloseFirestore() {
	if Client != nil {
		_ = Client.Close()
	}
}
