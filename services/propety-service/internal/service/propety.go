package service

import "github.com/HoanNghi16/rentix-services/propety-service/internal/repository"

type PropService struct {
	repo *repository.PropRepo
}

func NewService(repo *repository.PropRepo) *PropService {
	return &PropService{
		repo,
	}
}

func (service *PropService) GetComplexes(userID uint, role string) () {

}