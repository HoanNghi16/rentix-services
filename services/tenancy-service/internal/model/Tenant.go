package model

import (
	"time"

	"gorm.io/gorm"
)

type Tenant struct {
	UserID      uint       `gorm:"primaryKey;not null"` // use UserID
	Name        string     `gorm:"not null"`
	Email       string     `gorm:"not null"`
	Phone       string     `gorm:"not null"`
	DateOfBirth *time.Time `gorm:"default:null;type:DATE"`
	CitizenID string `gorm:"unique;not null"` //Số căn cước công dân
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}