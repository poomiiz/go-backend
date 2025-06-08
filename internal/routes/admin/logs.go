package routes

import (
	"context"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/iterator"
)

func RegisterLogsRoutes(r *gin.Engine, firestoreClient *firestore.Client) {
	r.GET("/admin/logs", func(c *gin.Context) {
		iter := firestoreClient.Collection("ai_logs").OrderBy("timestamp", firestore.Desc).Limit(20).Documents(context.Background())
		var logs []map[string]interface{}
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			data := doc.Data()
			logs = append(logs, data)
		}
		c.JSON(http.StatusOK, logs)
	})
}
