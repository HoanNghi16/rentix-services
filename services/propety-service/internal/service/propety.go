package service

import (
	"github.com/HoanNghi16/rentix-services/propety-service/internal/model"
	"github.com/HoanNghi16/rentix-services/propety-service/internal/repository"
)

type PropService struct {
	repo *repository.PropRepo
}

func NewService(repo *repository.PropRepo) *PropService {
	return &PropService{
		repo,
	}
}

func (service *PropService) GetComplexes(userID uint) ([]model.Complex, error) {
	return service.repo.FindComplexesByOwner(userID)
}


func (service *PropService) NewComplex(userID uint, input *model.Complex)(*model.Complex, error){
	input.OwnerID = userID
	if err := service.repo.CreateComplex(input); err != nil{
		return nil, err
	}
	return input, nil
}