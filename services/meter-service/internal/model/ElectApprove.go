package model

import "time"

type ElectApprove struct {
	ID         uint      `gorm:"primaryKey;autoIncrement"`
	SentByID   uint      `gorm:"not null;"`
	ApproveStatus ApproveStatus `gorm:"type:varchar(15);not null;default:PENDING"`
	UpdatedBy uint      `gorm:"not null;"`
	ImageURL   string    `gorm:"not null"`
	UploadedAt time.Time `gorm:"default:now()"`

	ReadingMeter float64 `gorm:"type:DECIMAL(7,1)"`
}

