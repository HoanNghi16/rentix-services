package model

import "time"


type EmployeeAssignment struct {
	EmployeeUserID    uint `gorm:"not null" json:"user_id"`
	ComplexID uint `gorm:"not null" json:"complex_id"`
	Complex *Complex `gorm:"foreignKey:ComplexID" json:"complex"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}