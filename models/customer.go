package models

type Customer struct {
	ID       uint   `gorm:"primaryKey"`
	Name     string `gorm:"size:100;not null"`
	Phone    string `gorm:"size:20;not null"`
	Address  string `gorm:"type:text;not null"`
	Category string `gorm:"type:enum('setia','reguler');default:'reguler'"`
}
