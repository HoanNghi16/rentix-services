package repository

import (
	"errors"
	"fmt"
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

func (repo *PropRepo) GetComplexes(userID uint, role string) ([]*grpcapi.Complex, error) {
	complexes := make([]*grpcapi.Complex, 0)
	switch role {
	case "owner":
		{
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
				return complexes, err
			}
			return complexes, nil
		}
	case "manager":
		{
			var complexes_ids []uint
			assign_err := repo.db.Model(&model.EmployeeAssignment{}).
				Select("complex_id").Where("employee_user_id = ? AND is_active = ?", userID, true).
				Scan(&complexes_ids).Error
			if assign_err != nil {
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
				return complexes, err
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

// CreateComplexWithOutbox commits the complex, employee assignments, and event
// together. The payload callback runs after the database assigns the complex ID.
func (repo *PropRepo) CreateComplexWithOutbox(complex *model.Complex, employeeUserIDs []uint, topic string, payloadForComplex func(uint) ([]byte, error)) error {
	return repo.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(complex).Error; err != nil {
			return err
		}

		for _, employeeUserID := range employeeUserIDs {
			assignment := model.EmployeeAssignment{
				EmployeeUserID: employeeUserID,
				ComplexID:      complex.ID,
				IsActive:       true,
			}
			if err := tx.Clauses(clause.OnConflict{
				Columns: []clause.Column{{Name: "employee_user_id"}, {Name: "complex_id"}},
				DoUpdates: clause.Assignments(map[string]interface{}{
					"is_active":  true,
					"updated_at": time.Now(),
				}),
			}).Create(&assignment).Error; err != nil {
				return err
			}
		}

		payload, err := payloadForComplex(complex.ID)
		if err != nil {
			return err
		}
		outbox := model.Outbox{
			Topic:         topic,
			Key:           fmt.Sprint(complex.ID),
			Payload:       payload,
			Status:        "pending",
			NextAttemptAt: time.Now(),
		}
		return tx.Create(&outbox).Error
	})
}

// ClaimOutboxBatch leases due events so multiple service instances can publish
// concurrently without normally claiming the same row.
func (repo *PropRepo) ClaimOutboxBatch(limit int, lease time.Duration) ([]model.Outbox, error) {
	var claimed []model.Outbox
	err := repo.db.Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE", Options: "SKIP LOCKED"}).
			Where("status = ? AND next_attempt_at <= ? OR status = ? AND lease_until <= ?", "pending", now, "processing", now).
			Order("id ASC").Limit(limit).Find(&claimed).Error; err != nil {
			return err
		}
		if len(claimed) == 0 {
			return nil
		}
		leaseUntil := now.Add(lease)
		ids := make([]uint, 0, len(claimed))
		for _, event := range claimed {
			ids = append(ids, event.ID)
		}
		if err := tx.Model(&model.Outbox{}).Where("id IN ?", ids).Updates(map[string]interface{}{
			"status": "processing", "lease_until": leaseUntil,
		}).Error; err != nil {
			return err
		}
		for i := range claimed {
			claimed[i].Status = "processing"
			claimed[i].LeaseUntil = &leaseUntil
		}
		return nil
	})
	return claimed, err
}

func (repo *PropRepo) MarkOutboxPublished(id uint) error {
	now := time.Now()
	return repo.db.Model(&model.Outbox{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status": "published", "published_at": now, "lease_until": nil, "last_error": "",
	}).Error
}

func (repo *PropRepo) MarkOutboxFailed(id uint, publishErr error) error {
	var event model.Outbox
	if err := repo.db.Select("attempts").First(&event, id).Error; err != nil {
		return err
	}
	attempts := event.Attempts + 1
	delay := time.Second * time.Duration(1<<min(attempts, 8))
	return repo.db.Model(&model.Outbox{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status": "pending", "attempts": attempts, "next_attempt_at": time.Now().Add(delay),
		"lease_until": nil, "last_error": publishErr.Error(),
	}).Error
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
			"is_active":  true,
			"updated_at": time.Now(),
		}),
	}).Create(&assignment).Error
}
