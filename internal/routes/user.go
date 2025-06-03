package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/poomiiz/go-backend/internal/services"
)

func RegisterUserRoutes(r *gin.Engine) {
	userSvc := services.NewUserService()
	grp := r.Group("/user")
	{
		grp.POST("/register", func(c *gin.Context) {
			var payload struct {
				Email    string `json:"email"`
				Password string `json:"password"`
				Role     string `json:"role"`
			}
			if err := c.ShouldBindJSON(&payload); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
				return
			}
			uid, err := userSvc.Register(c.Request.Context(), payload.Email, payload.Password, payload.Role)
			if err != nil {
				c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusCreated, gin.H{"userId": uid})
		})

		grp.POST("/login", func(c *gin.Context) {
			var payload struct {
				Email    string `json:"email"`
				Password string `json:"password"`
			}
			if err := c.ShouldBindJSON(&payload); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
				return
			}
			uid, user, err := userSvc.Login(c.Request.Context(), payload.Email, payload.Password)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"userId": uid, "role": user.Role})
		})

		grp.GET("/:id", func(c *gin.Context) {
			userID := c.Param("id")
			user, err := userSvc.GetByID(c.Request.Context(), userID)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
				return
			}
			c.JSON(http.StatusOK, user)
		})
	}
}
