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

	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware(db))

	api.POST("/logout", handlers.Logout(db))

	// Branch routes (admin only)
	admin := api.Group("/admin")
	admin.Use(middleware.AdminOnly())
	{
		// User routes
		admin.POST("/users", handlers.CreateUser(db))
		admin.GET("/users", handlers.GetUsers(db))
		admin.GET("/users/:id", handlers.GetUser(db))
		admin.PUT("/users/:id", handlers.UpdateUser(db))
		admin.DELETE("/users/:id", handlers.DeleteUser(db))

		// Branch routes
		admin.POST("/branches", handlers.CreateBranch(db))
		admin.GET("/branches", handlers.GetBranches(db))
		admin.GET("/branches/:id", handlers.GetBranch(db))
		admin.PUT("/branches/:id", handlers.UpdateBranch(db))
		admin.DELETE("/branches/:id", handlers.DeleteBranch(db))
	}
}
