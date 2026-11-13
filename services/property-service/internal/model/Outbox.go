package model

import "time"

type Outbox struct {
	ID            uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	Topic         string     `gorm:"not null;index" json:"topic"`
	Key           string     `gorm:"not null" json:"key"`
	Payload       []byte     `gorm:"type:jsonb;not null" json:"payload"`
	Status        string     `gorm:"not null;default:pending;index" json:"status"`
	Attempts      int        `gorm:"not null;default:0" json:"attempts"`
	NextAttemptAt time.Time  `gorm:"not null;index" json:"next_attempt_at"`
	LeaseUntil    *time.Time `gorm:"index" json:"lease_until,omitempty"`
	LastError     string     `json:"last_error,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	PublishedAt   *time.Time `json:"published_at,omitempty"`
}
