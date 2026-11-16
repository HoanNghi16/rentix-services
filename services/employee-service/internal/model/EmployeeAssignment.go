package model

import "gorm.io/datatypes"

type EmployeeAssignment struct {
	ComplexID   uint     `gorm:"primaryKey"`
	EmployeeID  uint     `gorm:"primaryKey"`
	Employee    Employee `gorm:"foreignKey:EmployeeID"`
	IsActive    bool `gorm:"default:true"`
	Permissions datatypes.JSON
	CurrentRoleID *uint `gorm:"default:null"`
	CurrentRole *CustomRole `gorm:"foreignKey:CurrentRoleID"`
	CurrentRoleName string `gorm:"default:Quản lý"`
}