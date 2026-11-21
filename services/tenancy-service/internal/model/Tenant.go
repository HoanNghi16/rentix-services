package model

import (
	"time"

	"gorm.io/gorm"
)

type Tenant struct {
	UserID      uint       `gorm:"primaryKey;not null"` // use UserID
	Name        string     `gorm:"not null"`
	Email       string     `gorm:"not null"`
	DateOfBirth *time.Time `gorm:"default:null"`
	CitizenID string `gorm:"unique"` //Số căn cước công dân
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}