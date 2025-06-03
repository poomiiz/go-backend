// cmd/main.go
package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/poomiiz/go-backend/internal/routes"
	"github.com/poomiiz/go-backend/internal/utils"
)

func main() {
	// โหลด .env
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found or failed to load", err)
	}

	// ตรวจ GOOGLE_APPLICATION_CREDENTIALS
	if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") == "" {
		log.Fatal("GOOGLE_APPLICATION_CREDENTIALS is not set")
	}
	// ตรวจ FIREBASE_PROJECT_ID
	if os.Getenv("FIREBASE_PROJECT_ID") == "" {
		log.Fatal("FIREBASE_PROJECT_ID is not set")
	}

	// Init Firestore
	utils.InitFirestore()
	defer utils.CloseFirestore()

	// อ่าน PORT
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// สร้าง Gin router
	r := gin.Default()
	routes.RegisterUserRoutes(r)
	routes.RegisterCoinRoutes(r)
	routes.RegisterPackageRoutes(r)
	routes.RegisterPaymentRoutes(r)
	routes.RegisterNotificationRoutes(r)
	routes.RegisterAIRoutes(r)
	routes.RegisterReviewRoutes(r)
	routes.RegisterRankRoutes(r)
	routes.RegisterBookingRoutes(r)
	// ลงทะเบียน LINE Webhook
	routes.RegisterLineWebhook(r)

	// รันเซิร์ฟเวอร์
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
