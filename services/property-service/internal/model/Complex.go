package model

import "time"

type Complex struct {
	ID        uint        `gorm:"primaryKey;autoIncrement" json:"id"`
	OwnerID   uint        `gorm:"not null;uniqueIndex:owner_cplxname_idx" json:"owner_id"`
	Name      string      `gorm:"not null;uniqueIndex:owner_cplxname_idx" json:"name"`
	Avatar    string      `json:"avatar"`
	Address   string      `json:"address"`
	Ward      string      `json:"ward"`
	City      string      `json:"city"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}