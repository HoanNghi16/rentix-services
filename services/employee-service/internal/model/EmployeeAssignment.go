package model

import "gorm.io/datatypes"

type EmployeeAssignment struct {
	ComplexID       uint     `gorm:"primaryKey"`
	EmployeeID      uint     `gorm:"primaryKey"`
	Employee        Employee `gorm:"foreignKey:EmployeeID;references:UserID"`
	IsActive        bool     `gorm:"default:true"`
	IsCustomed      bool     `gorm:"default:true"`
	Permissions     datatypes.JSON
	CurrentRoleID   *uint         `gorm:"default:null"`
	CurrentRole     *EmployeeRole `gorm:"foreignKey:CurrentRoleID"`
	CurrentRoleName string        `gorm:"default:Manager"`
}
