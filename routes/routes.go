package routes

import (
	"laundre/handlers"
	"laundre/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB) {
	// Example: Health check
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// Public routes
	r.POST("/login", handlers.Login(db))

	// Protected routes
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	// Add your protected routes here

	// Add more routes here
}
