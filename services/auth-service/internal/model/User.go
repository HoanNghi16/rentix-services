package model

import "time"

type User struct {
	ID      uint   `gorm:"primaryKey;autoIncrement"`
	Email    string `gorm:"unique;not null"`
	Password  string `gorm:"not null"`
	Role string `gorm:"default:owner"` // owner | employee | tenant
	TokenVersion uint `gorm:"default:1"`
	CreatedAt time.Time
	UpdatedAt time.Time
	LastLogin time.Time
}

