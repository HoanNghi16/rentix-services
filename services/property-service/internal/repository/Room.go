package repository

import (
	"errors"

	"github.com/HoanNghi16/rentix-services/property-service/internal/model"
	"github.com/HoanNghi16/rentix-services/property-service/internal/repository/dto"
)

// filter will be handled, not now :(
func (repo *PropRepo) GetRoomsForOwner(userID uint, complexID uint, filter any)([]dto.RoomListItem, error) {
	rooms := make([]dto.RoomListItem, 0)
	if err := repo.db.Model(&model.Room{}).
		Select(`Complex.id as complex_id,
			rooms.id as id, 
			room_number, 
			current_end_date, 
			price,
			status`).
			Joins("JOIN complexes as Complex ON Complex.id = rooms.complex_id").
			Where("Complex.owner_id = ? AND Complex.id = ?", userID, complexID).Scan(&rooms).Error; err != nil{
		return rooms, err
	}
	return rooms, nil
}

func (repo *PropRepo) GetRoomsForEmployee(userID uint, complexID uint, filter any)([]dto.RoomListItem, error){
	var assignment *model.EmployeeAssignment;
	assign_err := repo.db.Where("employee_user_id = ? AND complex_id = ?", userID, complexID).First(&assignment).Error
	if assign_err != nil {
		return []dto.RoomListItem{},assign_err
	}
	if (assignment.ComplexID == complexID && assignment.EmployeeUserID == userID){
		rooms := make([]dto.RoomListItem, 0)
		if err := repo.db.Model(&model.Room{}).
		Select(`complex_id, id, room_number, current_end_date, price, status`).
		Where("complex_id = ?", complexID).
		Scan(&rooms).Error; err != nil{
			return rooms, err
		}
		return rooms, nil
	}
	return nil, errors.New("ASSIGNMENT_NOT_FOUND")
} 