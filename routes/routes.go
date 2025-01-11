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

	shared := api.Group("/shared")
	shared.Use(middleware.AdminOrStaffBranch())
	{
		shared.POST("/transaction", handlers.CreateTransaction(db))       // Create transaction
		shared.GET("/transaction", handlers.GetTransactions(db))          // Get all transactions
		shared.GET("/transaction/:id", handlers.GetTransaction(db))       // Get a transaction by ID
		shared.PUT("/transaction/:id", handlers.UpdateTransaction(db))    // Update a transaction
		shared.DELETE("/transaction/:id", handlers.DeleteTransaction(db)) // Delete a transaction

		shared.POST("/orders", handlers.CreateOrder(db))       // Create order
		shared.GET("/orders", handlers.GetOrders(db))          // Get all orders
		shared.GET("/orders/:id", handlers.GetOrder(db))       // Get order by ID
		shared.PUT("/orders/:id", handlers.UpdateOrder(db))    // Update order
		shared.DELETE("/orders/:id", handlers.DeleteOrder(db)) // Delete order

		shared.POST("/customers", handlers.CreateCustomer(db))       // Create customer
		shared.GET("/customers", handlers.GetCustomers(db))          // Get all customers
		shared.GET("/customers/:id", handlers.GetCustomer(db))       // Get customer by ID
		shared.PUT("/customers/:id", handlers.UpdateCustomer(db))    // Update customer
		shared.DELETE("/customers/:id", handlers.DeleteCustomer(db)) // Delete customer
	}

}
