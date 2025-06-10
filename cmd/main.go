// cmd/main.go

package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/poomiiz/go-backend/internal/routes"
	adminroutes "github.com/poomiiz/go-backend/internal/routes/admin"
	"github.com/poomiiz/go-backend/internal/utils"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ .env not found:", err)
	}
	log.Println("GOOGLE_APPLICATION_CREDENTIALS =", os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
	log.Println("FIREBASE_PROJECT_ID =", os.Getenv("FIREBASE_PROJECT_ID"))

	if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") == "" {
		log.Fatal("❌ GOOGLE_APPLICATION_CREDENTIALS is not set")
	}
	if os.Getenv("FIREBASE_PROJECT_ID") == "" {
		log.Fatal("❌ FIREBASE_PROJECT_ID is not set")
	}

	if err := utils.InitFirestore(); err != nil {
		log.Fatalf("🔥 Firestore init failed: %v", err)
	}
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
	routes.RegisterAIRoutes(r)
	routes.RegisterReviewRoutes(r)
	routes.RegisterRankRoutes(r)
	routes.RegisterBookingRoutes(r)
	routes.RegisterLineWebhook(r)
	// เรียกใช้จริงจาก internal/routes/admin
	adminroutes.RegisterPromptRoutes(r, utils.GetFirestoreClient())
	adminroutes.RegisterConfigRoutes(r, utils.GetFirestoreClient())
	adminroutes.RegisterLogsRoutes(r, utils.GetFirestoreClient())
	adminroutes.RegisterDeckRoutes(r, utils.GetFirestoreClient())
	// Health check
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":    "ok",
			"service":   "go-backend",
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})

	log.Println("🚀 Server started at port", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
