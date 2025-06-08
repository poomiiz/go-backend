package routes

import (
	"context"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
)

func RegisterDeckRoutes(r *gin.Engine, client *firestore.Client) {
	r.GET("/admin/decks", func(c *gin.Context) {
		// existing GetDecks logic...
	})

	r.GET("/admin/decks/:deckId/cards", func(c *gin.Context) {
		deckId := c.Param("deckId")
		ctx := context.Background()
		snapIter := client.Collection("decks").Doc(deckId).Collection("cards").Documents(ctx)
		var cards []map[string]interface{}
		for {
			doc, err := snapIter.Next()
			if err != nil {
				break
			}
			cards = append(cards, doc.Data())
		}
		c.JSON(http.StatusOK, cards)
	})
}
