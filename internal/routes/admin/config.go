package routes

import (
	"context"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
)

// RegisterConfigRoutes ผูก route /admin/config
func RegisterConfigRoutes(r *gin.Engine, client *firestore.Client) {
	group := r.Group("/admin")
	group.GET("/config", func(c *gin.Context) {
		ctx := context.Background()
		docs, err := client.Collection("config").Documents(ctx).GetAll()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		list := make([]map[string]interface{}, 0, len(docs))
		for _, d := range docs {
			m := d.Data()
			m["id"] = d.Ref.ID
			list = append(list, m)
		}
		c.JSON(http.StatusOK, list)
	})
}
