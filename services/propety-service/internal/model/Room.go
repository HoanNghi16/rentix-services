package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type Room struct {
	ID         uint       `gorm:"primaryKey;autoIncrement;"`
	ComplexID  uint       `gorm:"not null;uniqueIndex:idx_room_cmplx"`
	Complex    Complex    `gorm:"foreignKey:ComplexID"`
	RoomNumber string     `gorm:"not null;uniqueIndex:idx_room_cmplx" json:"room_number"`
	Status     RoomStatus `gorm:"not null;default:EMPTY"`
	Price decimal.Decimal `gorm:"type:NUMERIC(16,3);not null;"`
	CreatedAt  time.Time	
	UpdatedAt time.Time
}