package repository

import (
	"github.com/HoanNghi16/rentix-services/tenancy-service/internal/model"
	"gorm.io/gorm"
)

type TenancyRepo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB)(*TenancyRepo){
	return &TenancyRepo{
		db,
	}
}



func (repo *TenancyRepo) GetTenants(ComplexID uint)([]model.Tenant, error){
	var tenants []model.Tenant

	if err := repo.db.Model(&model.Tenant{}).Select(
		"name, phone, email, date_of_birth, room_number",
		).Joins("JOIN Tenancy tnc on tenants.user_id = tnc.tenant_id").
		Where("tnc.complex_id = ? AND is_active = true", ComplexID).Scan(tenants).Error; err != nil{
		return nil, err
	}
	return tenants, nil
} 