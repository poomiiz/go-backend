package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/poomiiz/go-backend/internal/services"
)

func RegisterPackageRoutes(r *gin.Engine) {
	coinSvc := services.NewCoinService()
	pkgSvc := services.NewPackageService(coinSvc)

	grp := r.Group("/package")
	{
		grp.POST("/buy", func(c *gin.Context) {
			var payload struct {
				UserID    string `json:"userId"`
				PackageID string `json:"packageId"`
			}
			if err := c.ShouldBindJSON(&payload); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
				return
			}
			up, err := pkgSvc.BuyPackage(c.Request.Context(), payload.UserID, payload.PackageID)
			if err != nil {
				c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, up)
		})

		grp.GET("/check", func(c *gin.Context) {
			userID := c.Query("userId")
			active, err := pkgSvc.CheckUserPackage(c.Request.Context(), userID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"active": active})
		})
	}
}
