package service

import (
	"errors"
	"sync"

	"github.com/HoanNghi16/rentix-services/property-service/internal/grpcapi"
	"github.com/HoanNghi16/rentix-services/property-service/internal/model"
	"github.com/HoanNghi16/rentix-services/property-service/internal/repository"
)

type PropService struct {
	repo *repository.PropRepo
}

func NewService(repo *repository.PropRepo) *PropService {
	return &PropService{
		repo,
	}
}

func (service *PropService) GetComplexes(userID uint64, role string) ([]grpcapi.Complex, error) {
	complexes, err := service.repo.GetComplexes(uint(userID), role)
	if err != nil {
		return nil, err
	}

	result := make([]grpcapi.Complex, 0, len(complexes))
	for _, complex := range complexes {
		result = append(result, grpcapi.Complex{
			ID:            uint64(complex.ID),
			Name:          complex.Name,
			Avatar:        complex.Avatar,
			Address:       complex.Address,
			Ward:          complex.Ward,
			City:          complex.City,
			TotalRooms:    int32(complex.TotalRooms),
			InactiveRooms: int32(complex.InactiveRooms),
		})
	}
	return result, nil
}

func (service *PropService) CreateComplex(request *grpcapi.CreateComplexRequest) (*grpcapi.Complex, error) {
	newComplex := &model.Complex{
		Name:    request.Name,
		Avatar:  request.Avatar,
		OwnerID: uint(request.UserID),
		Address: request.Address,
		Ward:    "",
		City:    "", //Null and handle in the future
	}

	newComplex, err := service.repo.CreateComplex(newComplex)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup

	for _, employeeUserID := range request.EmployeeUserIds {
		employeeID := uint(employeeUserID)
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := service.AddEmployeeToComplex(request.Role, uint(request.UserID), newComplex.ID, employeeID)
			if err != nil {
				// Handle error, log it or retry
			}
		}()
	}

	// handle Kafka event và tạo outbox, topic: complex.created, gửi trong payload gửi complex_id, owner_id, employee_user_ids, services, còn các thông tin khác để ko cần gửi

	wg.Wait()

	return &grpcapi.Complex{
		ID:            uint64(newComplex.ID),
		Name:          newComplex.Name,
		Avatar:        newComplex.Avatar,
		Address:       newComplex.Address,
		Ward:          newComplex.Ward,
		City:          newComplex.City,
		TotalRooms:    0,
		InactiveRooms: 0,
	}, nil
}

func (service *PropService) AddEmployeeToComplex(role string, userID uint, complexID uint, employeeUserID uint) error {
	if role != "owner" {
		return errors.New("UNAUTHORIZED")
	}
	return service.repo.AddEmployeeToComplex(userID, complexID, employeeUserID)
}
