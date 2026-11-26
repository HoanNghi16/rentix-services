package grpcapi

import (
	"context"

	"github.com/HoanNghi16/rentix-services/employee-service/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type employeeServer struct{ service *service.EmpService }

func NewEmployeeServer(employeeService *service.EmpService) EmployeeServiceServer {
	return &employeeServer{service: employeeService}
}

func (server *employeeServer) GetEmployees(_ context.Context, request *GetEmployeesRequest) (*GetEmployeesResponse, error) {
	if request == nil || request.OwnerID == 0 || request.OwnerID > uint64(^uint(0)) {
		return nil, status.Error(codes.InvalidArgument, "owner ID is required")
	}
	employees, err := server.service.GetEmployees(uint(request.OwnerID), request.Role)
	if err != nil {
		if err == service.ErrForbidden {
			return nil, status.Error(codes.PermissionDenied, "only owners can access employees")
		}
		return nil, status.Error(codes.Internal, "could not load employees")
	}
	response := &GetEmployeesResponse{Employees: make([]*Employee, 0, len(employees))}
	for _, employee := range employees {
		item := &Employee{
			UserID: uint64(employee.UserID), FullName: employee.FullName,
			Phone: employee.Phone, Email: employee.Email,
			Assignments: make([]*EmployeeAssignment, 0, len(employee.Assignments)),
		}
		for _, assignment := range employee.Assignments {
			item.Assignments = append(item.Assignments, &EmployeeAssignment{
				ComplexID: uint64(assignment.ComplexID),
			})
		}
		response.Employees = append(response.Employees, item)
	}
	return response, nil
}
