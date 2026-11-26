package model

type Tenancy struct {
	ID uint `gorm:"primaryKey;autoIncrement"`

	TenantID uint   `gorm:"not null;uniqueIndex:idx_room_tenant"`
	Tenant   Tenant `gorm:"foreignKey:TenantID"`

	ComplexID   uint   `gorm:"not null"`
	RoomID      uint   `gorm:"not null;uniqueIndex:idx_room_tenant"`
	RoomNummber string `gorm:"not null"`

	SharedContractID uint     `gorm:"not null"`
	SharedContract   Contract `gorm:"foreignKey:SharedContractID"` //Hợp đồng dùng chung cho nhiều khách cùng thuê 1 phòng

	IsActive bool `gorm:"default:true"`
}