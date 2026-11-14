package service

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"math"
	"time"

	"github.com/HoanNghi16/rentix-services/property-service/internal/grpcapi"
	"github.com/HoanNghi16/rentix-services/property-service/internal/model"
	"github.com/HoanNghi16/rentix-services/property-service/internal/repository"
	"github.com/HoanNghi16/rentix-services/property-service/internal/repository/dto"
)

type PropService struct {
	repo *repository.PropRepo
}

type EventPublisher interface {
	Publish(context.Context, string, string, []byte) error
}

const ComplexCreatedTopic = "complex.created"

type ComplexCreatedEvent struct {
	EventType       string            `json:"event_type"`
	ComplexID       uint              `json:"complex_id"`
	OwnerID         uint              `json:"owner_id"`
	EmployeeUserIDs []uint            `json:"employee_user_ids"`
	Services        *grpcapi.Services `json:"services"`
	OccurredAt      time.Time         `json:"occurred_at"`
}

func NewService(repo *repository.PropRepo) *PropService {
	return &PropService{
		repo,
	}
}

func (service *PropService) GetComplexes(userID uint64, role string) (*grpcapi.GetComplexesResponse, error) {
	log.Print("Đã tới đây")
	complexes, err := service.repo.GetComplexes(uint(userID), role)
	if err != nil {
		return nil, err
	}
	return &grpcapi.GetComplexesResponse{Complexes: complexes}, nil
}

func (service *PropService) GetRooms(request *grpcapi.GetRoomsRequest) (*grpcapi.GetRoomsResponse, error) {
	if request == nil || request.ComplexID <= 0 || request.UserID <= 0 || request.Role == "" {
		return nil, errors.New("INVALID_ARGUMENT")
	}

	userID := uint(request.UserID)
	complexID := uint(request.ComplexID)
	var rooms []dto.RoomListItem
	var err error
	switch request.Role {
	case "owner":
		rooms, err = service.repo.GetRoomsForOwner(userID, complexID, nil)
	case "manager":
		rooms, err = service.repo.GetRoomsForEmployee(userID, complexID, nil)
	default:
		return nil, errors.New("INVALID_ROLE")
	}
	if err != nil {
		return nil, err
	}

	response := &grpcapi.GetRoomsResponse{Rooms: make([]*grpcapi.Room, 0, len(rooms))}
	for _, room := range rooms {
		currentEndDate := ""
		if !room.CurrentEndDate.IsZero() {
			currentEndDate = room.CurrentEndDate.Format(time.RFC3339)
		}
		price, _ := room.Price.Float64()
		if math.IsInf(price, 0) || math.IsNaN(price) {
			return nil, errors.New("INVALID_ROOM_PRICE")
		}
		response.Rooms = append(response.Rooms, &grpcapi.Room{
			ComplexID:      int64(room.ComplexID),
			ID:             int64(room.ID),
			CurrentEndDate: currentEndDate,
			RoomNumber:     room.RoomNumber,
			Status:         string(room.Status),
			Price:          float32(price),
		})
	}
	return response, nil
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

	if request.Role != "owner" {
		return nil, errors.New("UNAUTHORIZED")
	}
	employeeUserIDs := make([]uint, 0, len(request.EmployeeUserIds))
	for _, employeeUserID := range request.EmployeeUserIds {
		employeeUserIDs = append(employeeUserIDs, uint(employeeUserID))
	}
	err := service.repo.CreateComplexWithOutbox(newComplex, employeeUserIDs, ComplexCreatedTopic, func(complexID uint) ([]byte, error) {
		return json.Marshal(ComplexCreatedEvent{
			EventType:       ComplexCreatedTopic,
			ComplexID:       complexID,
			OwnerID:         uint(request.UserID),
			EmployeeUserIDs: employeeUserIDs,
			Services:        request.Services,
			OccurredAt:      time.Now().UTC(),
		})
	})
	if err != nil {
		return nil, err
	}

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

func (service *PropService) PublishOutbox(ctx context.Context, publisher EventPublisher) error {
	events, err := service.repo.ClaimOutboxBatch(50, time.Minute)
	if err != nil {
		return err
	}
	var publishErr error
	for _, event := range events {
		if err := publisher.Publish(ctx, event.Topic, event.Key, event.Payload); err != nil {
			log.Printf("publish outbox event %d to %s: %v", event.ID, event.Topic, err)
			if markErr := service.repo.MarkOutboxFailed(event.ID, err); markErr != nil {
				log.Printf("mark outbox event %d for retry: %v", event.ID, markErr)
			}
			publishErr = errors.Join(publishErr, err)
			continue
		}
		if err := service.repo.MarkOutboxPublished(event.ID); err != nil {
			log.Printf("mark outbox event %d published: %v", event.ID, err)
			publishErr = errors.Join(publishErr, err)
		}
	}
	return publishErr
}

func (service *PropService) AddEmployeeToComplex(role string, userID uint, complexID uint, employeeUserID uint) error {
	if role != "owner" {
		return errors.New("UNAUTHORIZED")
	}
	return service.repo.AddEmployeeToComplex(userID, complexID, employeeUserID)
}
