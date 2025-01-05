package models

import "time"

type Transaction struct {
	ID            uint      `gorm:"primaryKey"`
	BranchID      uint      `gorm:"not null"`
	OrderID       uint      `gorm:"not null"`
	UserID        uint      `gorm:"not null"`
	TotalPrice    float64   `gorm:"type:decimal(10,2);not null"`
	PaymentStatus string    `gorm:"type:enum('paid','unpaid');default:'unpaid'"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	Branch        Branch    `gorm:"constraint:OnDelete:CASCADE"`
	Order         Order     `gorm:"constraint:OnDelete:CASCADE"`
	User          User      `gorm:"constraint:OnDelete:CASCADE"`
}
