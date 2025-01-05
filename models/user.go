package models

import "time"

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Username  string `gorm:"size:50;unique;not null"`
	Password  string `gorm:"size:255;not null"`
	Role      string `gorm:"type:enum('admin','staf');not null"`
	BranchID  *uint  `gorm:"default:null"`
	Status    string `gorm:"type:enum('active','inactive');default:'active'"`
	LastLogin *time.Time
}
