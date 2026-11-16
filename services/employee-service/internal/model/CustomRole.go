package model

import (
	"time"

	"gorm.io/datatypes"
)

type CustomRole struct {
	ID          uint `gorm:"primaryKey;autoIncrement"`
	OwnerID     uint `gorm:"uniqueIndex:idx_role_owner"`
	RoleName    string `gorm:"uniqueIndex:idx_role_owner"`
	Permissions datatypes.JSON
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
