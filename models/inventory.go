package models

type Inventory struct {
	ID       uint   `gorm:"primaryKey"`
	BranchID uint   `gorm:"not null"`
	Name     string `gorm:"size:100;not null"`
	Stock    int    `gorm:"not null"`
	Branch   Branch `gorm:"constraint:OnDelete:CASCADE"`
}
