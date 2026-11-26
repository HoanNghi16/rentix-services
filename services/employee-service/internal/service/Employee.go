package service

import (
	"errors"

	"github.com/HoanNghi16/rentix-services/employee-service/internal/model"
	"github.com/HoanNghi16/rentix-services/employee-service/internal/repository"
)

var ErrForbidden = errors.New("FORBIDDEN")

type EmpService struct{ repo *repository.EmpRepo }

func NewService(repo *repository.EmpRepo) *EmpService { return &EmpService{repo: repo} }

func (service *EmpService) GetEmployees(ownerID uint, role string) ([]model.Employee, error) {
	if role != "owner" {
		return nil, ErrForbidden
	}
	return service.repo.GetEmployees(ownerID)
}
