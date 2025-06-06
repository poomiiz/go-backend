package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/poomiiz/go-backend/internal/routes" // import routes อย่างเดียว
	"github.com/poomiiz/go-backend/internal/utils"
)

func main() {
	// โหลด .env
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env not found", err)
	}

	if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") == "" {
		log.Fatal("GOOGLE_APPLICATION_CREDENTIALS is not set")
	}
	if os.Getenv("FIREBASE_PROJECT_ID") == "" {
		log.Fatal("FIREBASE_PROJECT_ID is not set")
	}

	utils.InitFirestore()
	defer utils.CloseFirestore()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := gin.Default()
	routes.RegisterUserRoutes(r)
	routes.RegisterCoinRoutes(r)
	routes.RegisterPackageRoutes(r)
	routes.RegisterPaymentRoutes(r)
	routes.RegisterNotificationRoutes(r)
	routes.RegisterAIRoutes(r) // ลงทะเบียน ai routes
	routes.RegisterReviewRoutes(r)
	routes.RegisterRankRoutes(r)
	routes.RegisterBookingRoutes(r)
	routes.RegisterLineWebhook(r)

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
