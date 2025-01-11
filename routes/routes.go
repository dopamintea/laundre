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
	branches := api.Group("/branches")
	branches.Use(middleware.AdminOnly())
	{
		branches.POST("/", handlers.CreateBranch(db))
		branches.GET("/", handlers.GetBranches(db))
		branches.GET("/:id", handlers.GetBranch(db))
		branches.PUT("/:id", handlers.UpdateBranch(db))
		branches.DELETE("/:id", handlers.DeleteBranch(db))
	}
}
