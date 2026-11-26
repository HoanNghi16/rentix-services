package grpcapi

import (
	"context"

	"google.golang.org/grpc"
)

const employeeServiceName = "rentix.employee.v1.EmployeeService"

type EmployeeServiceServer interface {
	GetEmployees(context.Context, *GetEmployeesRequest) (*GetEmployeesResponse, error)
}

func RegisterEmployeeServiceServer(registrar grpc.ServiceRegistrar, server EmployeeServiceServer) {
	registrar.RegisterService(&EmployeeService_ServiceDesc, server)
}

func getEmployeesHandler(server any, ctx context.Context, decode func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	request := new(GetEmployeesRequest)
	if err := decode(request); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return server.(EmployeeServiceServer).GetEmployees(ctx, request)
	}
	info := &grpc.UnaryServerInfo{Server: server, FullMethod: "/" + employeeServiceName + "/GetEmployees"}
	handler := func(ctx context.Context, req any) (any, error) {
		return server.(EmployeeServiceServer).GetEmployees(ctx, req.(*GetEmployeesRequest))
	}
	return interceptor(ctx, request, info, handler)
}

var EmployeeService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: employeeServiceName,
	HandlerType: (*EmployeeServiceServer)(nil),
	Methods:     []grpc.MethodDesc{{MethodName: "GetEmployees", Handler: getEmployeesHandler}},
	Streams:     []grpc.StreamDesc{},
	Metadata:    "employee.proto",
}
