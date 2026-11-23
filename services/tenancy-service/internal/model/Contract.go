package model

import "time"

type Contract struct {
	ID        uint      `gorm:"not null;primaryKey;autoIncrement"`
	OwnerID   uint      `gorm:"not null"`
	TenantID  uint      `gorm:"not null"`
	Tenant    *Tenant   `gorm:"foreignKey:TenantID"`

	RoomID uint `gorm:"not null"`

	StartDate time.Time `gorm:"type:DATE;not null"`
	EndDate time.Time `gorm:"type:DATE;not null"`

}
