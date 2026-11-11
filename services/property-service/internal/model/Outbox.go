package model

type Outbox struct {
	ID uint `gorm:"primaryKey;autoIncrement" json:"id"`
	
}