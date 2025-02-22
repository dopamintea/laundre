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

		admin.GET("/finance/profit", handlers.GetProfit(db))
		admin.GET("/finance/profit/:branch_id", handlers.GetProfitByBranch(db))
		admin.GET("/finance/gross", handlers.GetGrossProfit(db))
		admin.GET("/finance/gross/:branch_id", handlers.GetGrossProfitByBranch(db))
		admin.POST("/transaction/report", handlers.GetTransactionByDate(db))
		admin.GET("/transaction/report/:branch_id", handlers.GetTransactionsByBranch(db))
		admin.GET("/expense/branch/:branch_id", handlers.GetExpensesByBranch(db))
	}

	shared := api.Group("/shared")
	shared.Use(middleware.AdminOrStaffBranch())
	{
		shared.POST("/transaction", handlers.CreateTransaction(db))
		shared.GET("/transaction", handlers.GetTransactions(db))
		shared.GET("/transaction/:id", handlers.GetTransaction(db))
		shared.PUT("/transaction/:id", handlers.UpdateTransaction(db))
		shared.GET("/transaction/status/:status", handlers.GetTransactionsByOrderStatus(db))
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

		shared.POST("/inventory", handlers.CreateInventory(db))
		shared.GET("/inventory", handlers.GetAllInventories(db))
		shared.GET("/inventory/:id", handlers.GetInventoryByID(db))
		shared.PUT("/inventory/:id", handlers.UpdateInventory(db))
		shared.DELETE("/inventory/:id", handlers.DeleteInventory(db))
		shared.POST("/inventory/branch", handlers.GetInventoryByBranch(db))

		shared.POST("/expense", handlers.CreateExpense(db))
		shared.GET("/expense", handlers.GetAllExpenses(db))
		shared.GET("/expense/:id", handlers.GetExpenseByID(db))
		shared.PUT("/expense/:id", handlers.UpdateExpense(db))
		shared.DELETE("/expense/:id", handlers.DeleteExpense(db))
	}

}
