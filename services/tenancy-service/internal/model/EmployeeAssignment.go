package model

type EmployeeAssignment struct {
	ID             uint `gorm:"primaryKey;autoIncrement"`
	OwnerID        uint `gorm:"not null;uniqueIndex:idx_owner_assign"`
	ComplexID      uint `gorm:"not null;uniqueIndex:idx_owner_assign"`
	EmployeeUserID uint `gorm:"not null;uniqueIndex:idx_owner_assign"`
	IsActive       bool `gorm:"default:true"`

	PermissionVersion int32 `gorm:"default:0"`

	// tables' name. ex: ("contracts", "tenants",...)
	ResourceName string `gorm:"not null;uniqueIndex:idx_owner_assign"`

	// level permission
	ReadAllowed   bool `gorm:"default:true"`
	CreateAllowed bool `gorm:"default:true"`
	UpdateAllowed bool `gorm:"default:true"`
}