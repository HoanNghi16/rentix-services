package model

import "time"


type EmployeeAssignment struct {
	ID uint `gorm:"primaryKey;autoIncrement" json:"id"`
	EmployeeUserID    uint `gorm:"not null;uniqueIndex:idx_employee_complex" json:"user_id"`
	ComplexID uint `gorm:"not null;uniqueIndex:idx_employee_complex" json:"complex_id"`
	Complex *Complex `gorm:"foreignKey:ComplexID" json:"complex"`
	IsActive bool `gorm:"not null; default:true" json:"is_active"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}