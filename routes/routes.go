package routes

import (
	"laundre/handlers"
	"laundre/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB) {

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	r.POST("/login", handlers.Login(db))

	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware(db))

	protected.POST("/logout", handlers.Logout(db))

}
