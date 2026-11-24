package model

type BillingRoom struct {
	ID          uint `gorm:"primaryKey;autoIncrement"`
	OwnerID     uint `gorm:"not null"`
	ComplexID   uint `gorm:"uniqueIndex:idx_cmplx_billing_room"`
	RoomID      uint `gorm:"uniqueIndex:idx_cmplx_billing_room"`
	TenantCount int  `gorm:"default:0"`
	IsActive    bool `gorm:"default:false"`
}