package model

import (
	"github.com/shopspring/decimal"
)

type RoomResponse struct {
	ID         uint            `json:"id"`
	ComplexID  uint            `json:"complex_id"`
	RoomNumber string          `json:"room_number"`
	Status     RoomStatus      `json:"status"`
	Price      decimal.Decimal `json:"price"`
}