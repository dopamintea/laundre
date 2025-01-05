package models

import "time"

type Expense struct {
	ID          uint      `gorm:"primaryKey"`
	BranchID    uint      `gorm:"not null"`
	Description string    `gorm:"type:text;not null"`
	Amount      float64   `gorm:"type:decimal(10,2);not null"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	Branch      Branch    `gorm:"constraint:OnDelete:CASCADE"`
}
