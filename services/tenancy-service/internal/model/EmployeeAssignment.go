package model

type EmployeeAssignment struct {
	ID             uint   `gorm:"primaryKey;autoIncrement"`
	OwnerID        uint   `gorm:"not null;uniqueIndex:idx_owner_assign"`
	ComplexID      uint   `gorm:"not null;uniqueIndex:idx_owner_assign"`
	EmployeeUserID uint   `gorm:"not null;uniqueIndex:idx_owner_assign"`
	IsActive       bool   `gorm:"default:true"`
	ResourceName   string `gorm:"not null;uniqueIndex:idx_owner_assign"` //Tên bảng vd: ("contracts", "tenants",...)
	ReadAllowed    bool   `gorm:"default:true"`
	CreateAllowed  bool   `gorm:"default:true"`
	UpdateAllowed  bool   `gorm:"default:true"`
}