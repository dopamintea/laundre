package models

import "time"

type Order struct {
	ID         uint      `gorm:"primaryKey"`
	BranchID   uint      `gorm:"not null"`
	CustomerID uint      `gorm:"not null"`
	Status     string    `gorm:"type:enum('masuk','proses','urgent', 'done');default:'masuk'"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
	Price      float64   `gorm:"type:decimal(10,2);not null"`
	Branch     Branch    `gorm:"constraint:OnDelete:CASCADE"`
	Customer   Customer  `gorm:"constraint:OnDelete:CASCADE"`
}
