package repository

import (
	"github.com/HoanNghi16/rentix-services/employee-service/internal/model"
	"gorm.io/gorm"
)

type EmpRepo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *EmpRepo { return &EmpRepo{db: db} }

func (repo *EmpRepo) GetEmployees(ownerID uint) ([]model.Employee, error) {
	var employees []model.Employee
	err := repo.db.Model(&model.Employee{}).
		Select("employees.user_id", "employees.full_name", "employees.email", "employees.phone").
		Where("employees.owner_id = ? AND employees.is_active = ?", ownerID, true).
		Preload("Assignments", func(db *gorm.DB) *gorm.DB {
			return db.Select("employee_id", "complex_id").Where("is_active = ?", true)
		}).
		Find(&employees).Error
	return employees, err
}

func (repo *EmpRepo) CreateEmployee (emp model.Employee, complex_ids []uint) (*model.Employee, error){
	return &emp, repo.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(&emp).Error
		if err != nil {
			return err
		}
		if(len(complex_ids) == 0){
			return nil
		}

		assignments := make([]model.EmployeeAssignment, 0, len(complex_ids))
		for _, id := range complex_ids{
			assignments = append(assignments, model.EmployeeAssignment{
				ComplexID: id,
				EmployeeID: emp.UserID,
				IsActive: true,
				CurrentRoleID: nil,
				CurrentRoleName: "Quản lý",
				IsCustomed: false,
				Permissions: nil,

			})
		}
		return tx.Create(&assignments).Error
	})
}
