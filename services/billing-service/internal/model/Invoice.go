package model

type Invoice struct {
	ID           uint          `gorm:"primaryKey;autoIncrement"`
	RoomID       uint          `gorm:"uniqueIndex:idx_bill_room_per_month"`
	BillingYear  int           `gorm:"uniqueIndex:idx_bill_room_per_month"`
	BillingMonth int           `gorm:"uniqueIndex:idx_bill_room_per_month"`
	TotalAmount  float64       `gorm:"type:DECIMAL(15,3)"`
	Status       InvoiceStatus `gorm:"not null;default:PENDING"`
	Note         *string       // ghi chú thêm nếu có

	Discount   float32 `gorm:"max:1;min:0;type:DECIMAL(15,3)"`
	DiscountBy string  `gorm:"default:percent"` // percent | VNĐ

	Items []InvoiceItem
}

type InvoiceItem struct {
	ID          uint    `gorm:"primaryKey;autoIncrement"`
	InvoiceID   uint    `gorm:"uniqueIndex:idx_service_invoice"`
	Invoice     Invoice `gorm:"foreignKey:InvoiceID"`
	ServiceName string  `gorm:"uniqueIndex:idx_service_invoice"`
	Unit        Unit
	Quantity    float64 `gorm:"type:DECIMAL(7,1)"`
	UnitPrice   float64 `gorm:"type:DECIMAL(15,3)"`
	SubTotal    float64 `gorm:"type:DECIMAL(15,3)"`
}