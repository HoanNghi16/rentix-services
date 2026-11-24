package model

type RoomServicePrice struct {
	ID          uint    `gorm:"primaryKey;autoIncrement"`
	RoomID      uint    `gorm:"not null;uniqeIndex:idx_room_service"`
	ServiceName string  `gorm:"uniqeIndex:idx_room_service"`
	UnitPrice   float64 `gorm:"type:DECIMAL(15,3)"`
	Unit        Unit    `gorm:"type:VARCHAR(20)"`
}
