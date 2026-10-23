package repository

import (
	"github.com/HoanNghi16/rentix-services/propety-service/internal/model"
	"gorm.io/gorm"
)

type PropRepo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *PropRepo {
	return &PropRepo{
		db,
	}
}

func (repo *PropRepo) CreateComplex(input *model.Complex) error {
	return repo.db.Create(input).Error
}

func (repo *PropRepo) CreateRoom(input *model.Room) error{
	return repo.db.Create(input).Error
}

func (repo *PropRepo) FindRoomsByComplex(complexID uint) ([]model.Room, error){
	var rooms []model.Room;

	if sl_err := repo.db.Where("complex_id=?", complexID).Find(&rooms).Error; sl_err != nil{
		return make([]model.Room, 0), sl_err
	}
	return rooms, nil
}

func (repo *PropRepo) FindRoomByID(roomID uint) (*model.Room, error){
	var room model.Room;
	if err := repo.db.Where("id=?", roomID).First(&room).Error; err != nil{
		return nil, err
	}
	return &room, nil
}


func (repo *PropRepo) FindComplexByOwner(ownerID uint) (*model.Complex, error) {
	var complex *model.Complex
	if err := repo.db.Where("owner_id=?", ownerID).First(&complex).Error; err != nil{
		return nil, err
	}
	return complex, nil
}

func (repo *PropRepo) FindComplexesByOwner(ownerID uint) ([]model.Complex, error){
	var complexes []model.Complex
	if err := repo.db.Where("owner_id=?", ownerID).Find(&complexes).Error; err != nil{
		return nil, err
	}
	return complexes, nil
}
