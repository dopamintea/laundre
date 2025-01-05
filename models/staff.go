package models

type Staff struct {
	ID       uint   `gorm:"primaryKey"`
	Name     string `gorm:"size:100;not null"`
	BranchID uint   `gorm:"not null"`
	Phone    string `gorm:"size:20;not null"`
	Branch   Branch `gorm:"constraint:OnDelete:CASCADE"`
}
