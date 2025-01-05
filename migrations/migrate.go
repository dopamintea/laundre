package migrations

import (
	"laundre/models"
	"log"

	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) {
	err := db.AutoMigrate(
		&models.User{},
		&models.Branch{},
		&models.Inventory{},
		&models.Customer{},
		&models.Order{},
		&models.Transaction{},
		&models.Expense{},
		&models.Staff{},
		&models.Log{},
	)
	if err != nil {
		log.Fatal("Failed to run migrations:", err)
	} else {
		log.Println("Database migrated successfully!")
	}
}
