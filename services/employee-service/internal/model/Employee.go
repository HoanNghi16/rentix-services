package model

import "time"

type Employee struct {
	UserID      uint `gorm:"primaryKey;autoIncrement:false;not null"` //Use UserID for EmployeeID
	OwnerID     *uint
	Avatar      string
	FullName    string
	DateOfBirth *time.Time `gorm:"default:null"`

	PermissionVersion int32 `gorm:"default:0"`

	// Lấy danh sách phân công/complex_ids
	Assignments []EmployeeAssignment `gorm:"foreignKey:EmployeeID;references:UserID"`

	Phone     string
	Email     string
	Gender    string
	CreatedAt time.Time
	UpdatedAt time.Time
	IsActive  bool `gorm:"default:false"`
}
