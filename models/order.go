package models

import "time"

type Order struct {
	ID         uint      `gorm:"primaryKey"`
	BranchID   uint      `gorm:"not null"`
	CustomerID uint      `gorm:"not null"`
	Status     string    `gorm:"type:enum('masuk','proses','urgent');default:'masuk'"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
	Branch     Branch    `gorm:"constraint:OnDelete:CASCADE"`
	Customer   Customer  `gorm:"constraint:OnDelete:CASCADE"`
}
