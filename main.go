package main

import (
	"laundre/config"
	"laundre/migrations"
	"laundre/routes"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Connect to the database
	db, err := config.ConnectDatabase()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Run migrations
	migrations.RunMigrations(db)

	// Initialize Gin router
	r := gin.Default()

	// Register routes
	routes.RegisterRoutes(r, db)

	// Start the server
	log.Println("Server is running on http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
