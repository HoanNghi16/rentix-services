package dto

import (
	"time"

	"github.com/HoanNghi16/rentix-services/property-service/internal/model"
	"github.com/shopspring/decimal"
)

type RoomListItem struct {
	ID         uint           
	ComplexID  uint           
	RoomNumber string         
	CurrentEndDate time.Time	
	Status     model.RoomStatus     
	Price      decimal.Decimal
}