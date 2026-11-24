package model

import "time"

// Only Owner can modify
type ComplexServicePrice struct {
	ID  uint `gorm:"primaryKey;autoIncrement"`
	ComplexID uint `gorm:"not null;uniqueIndex:idx_cplx_srvce"`
	OwnerID uint `gorm:"not null"`
	UpdatedAt time.Time
	CreatedAt time.Time
	ServiceName string `gorm:"uniqueIndex:idx_cplx_srvce"`
	UnitPrice float64 `gorm:"type:DECIMAL(15,3)"`
	Unit Unit `gorm:"type:VARCHAR(20)"`
}