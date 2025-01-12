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

	admin := api.Group("/admin")
	admin.Use(middleware.AdminOnly())
	{
		admin.POST("/users", handlers.CreateUser(db))
		admin.GET("/users", handlers.GetUsers(db))
		admin.GET("/users/:id", handlers.GetUser(db))
		admin.PUT("/users/:id", handlers.UpdateUser(db))
		admin.DELETE("/users/:id", handlers.DeleteUser(db))

		admin.POST("/branches", handlers.CreateBranch(db))
		admin.GET("/branches", handlers.GetBranches(db))
		admin.GET("/branches/:id", handlers.GetBranch(db))
		admin.PUT("/branches/:id", handlers.UpdateBranch(db))
		admin.DELETE("/branches/:id", handlers.DeleteBranch(db))
	}

	shared := api.Group("/shared")
	shared.Use(middleware.AdminOrStaffBranch())
	{
		shared.POST("/transaction", handlers.CreateTransaction(db))
		shared.GET("/transaction", handlers.GetTransactions(db))
		shared.GET("/transaction/:id", handlers.GetTransaction(db))
		shared.PUT("/transaction/:id", handlers.UpdateTransaction(db))
		shared.DELETE("/transaction/:id", handlers.DeleteTransaction(db))

		shared.POST("/orders", handlers.CreateOrder(db))
		shared.GET("/orders", handlers.GetOrders(db))
		shared.GET("/orders/:id", handlers.GetOrder(db))
		shared.PUT("/orders/:id", handlers.UpdateOrder(db))
		shared.DELETE("/orders/:id", handlers.DeleteOrder(db))

		shared.POST("/customers", handlers.CreateCustomer(db))
		shared.GET("/customers", handlers.GetCustomers(db))
		shared.GET("/customers/:id", handlers.GetCustomer(db))
		shared.PUT("/customers/:id", handlers.UpdateCustomer(db))
		shared.DELETE("/customers/:id", handlers.DeleteCustomer(db))
	}

}
