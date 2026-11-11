package repository

import (
	"errors"
	"time"

	"github.com/HoanNghi16/rentix-services/property-service/internal/grpcapi"
	"github.com/HoanNghi16/rentix-services/property-service/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PropRepo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *PropRepo {
	return &PropRepo{
		db,
	}
}

func (repo *PropRepo) GetComplexes(userID uint, role string)([]grpcapi.Complex, error){
	var complexes []grpcapi.Complex;
	switch role{
		case "owner":{
			err := repo.db.
				Model(&model.Complex{}).
				Select(`
					complexes.*,
					COUNT(rooms.id) AS total_rooms,
					COUNT(rooms.id) FILTER (WHERE rooms.status = ?) AS inactive_rooms
				`, "EMPTY").
				Joins("LEFT JOIN rooms ON rooms.complex_id = complexes.id").
				Where("complexes.owner_id = ?", userID).
				Group("complexes.id").
				Scan(&complexes).Error

			if err != nil {
				return nil, err
			}
			return complexes, nil
		}
		case "manager":{
			var complexes_ids []uint
			assign_err := repo.db.Model(&model.EmployeeAssignment{}).
				Select("complex_id").Where("employee_user_id = ? AND is_active = ?", userID, true).
				Scan(&complexes_ids).Error;
			if assign_err != nil{
				return nil, assign_err
			}  
			err := repo.db.
				Model(&model.Complex{}).
				Select(`
					complexes.*,
					COUNT(rooms.id) AS total_rooms,
					COUNT(rooms.id) FILTER (WHERE rooms.status = ?) AS inactive_rooms
				`, "EMPTY").
				Joins("LEFT JOIN rooms ON rooms.complex_id = complexes.id").
				Where("complexes.id IN ?", complexes_ids).
				Group("complexes.id").
				Scan(&complexes).Error
			if err != nil {
				return nil, err
			}
			return complexes, nil
		}
	}
	return nil, errors.New("INVALID_ROLE")
}


func (repo *PropRepo) CreateComplex(complex *model.Complex) (*model.Complex, error) {
    if err := repo.db.Create(complex).Error; err != nil {
        return nil, err
    }
    return complex, nil
}



func (repo *PropRepo) AddEmployeeToComplex(userID uint, complexID uint, employeeUserID uint) error {
    assignment := model.EmployeeAssignment{
        EmployeeUserID: employeeUserID,
        ComplexID:      complexID,
        IsActive:       true,
    }

    return repo.db.Clauses(clause.OnConflict{
        Columns: []clause.Column{
            {Name: "employee_user_id"},
            {Name: "complex_id"},
        },
        DoUpdates: clause.Assignments(map[string]interface{}{
            "is_active": true,
            "updated_at": time.Now(),
        }),
    }).Create(&assignment).Error
}