package models

type Branch struct {
	ID      uint   `gorm:"primaryKey"`
	Name    string `gorm:"size:100;not null"`
	Address string `gorm:"type:text;not null"`
	Phone   string `gorm:"size:20;not null"`
}
