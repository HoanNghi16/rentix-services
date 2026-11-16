package model

import "time"

type Employee struct {
	ID          uint `gorm:"primaryKey;autoIncrement:false"` //Use UserID for EmployeeID
	OwnerID     uint
	Avatar string
	FullName    string
	DateOfBirth *time.Time `gorm:"default:null"`
	Phone	string
	Email string
	Gender bool
	CreatedAt time.Time
	UpdatedAt time.Time
}