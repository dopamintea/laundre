package models

import "time"

type Log struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null"`
	LoginTime time.Time `gorm:"autoCreateTime"`
	IPAddress string    `gorm:"size:45;not null"`
	User      User      `gorm:"constraint:OnDelete:CASCADE"`
}
