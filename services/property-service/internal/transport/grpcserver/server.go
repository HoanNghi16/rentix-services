package grpcserver

import (
	"context"
	"strings"

	"github.com/HoanNghi16/rentix-services/property-service/internal/grpcapi"
	"github.com/HoanNghi16/rentix-services/property-service/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const propertyServiceName = "rentix.property.v1.PropertyService"

type PropertyServiceServer interface {
	GetComplexes(context.Context, *grpcapi.GetComplexesRequest) (*grpcapi.GetComplexesResponse, error)
	CreateComplex(context.Context, *grpcapi.CreateComplexRequest) (*grpcapi.CreateComplexResponse, error)
}

type propertyServer struct{ service *service.PropService }

func NewPropertyServer(propertyService *service.PropService) PropertyServiceServer {
	return &propertyServer{service: propertyService}
}

func (server *propertyServer) GetComplexes(_ context.Context, request *grpcapi.GetComplexesRequest) (*grpcapi.GetComplexesResponse, error) {
	if request == nil || request.UserID == 0 || strings.TrimSpace(request.Role) == "" {
		return nil, status.Error(codes.InvalidArgument, "authenticated user claims are required")
	}
	complexes, err := server.service.GetComplexes(request.UserID, request.Role)
	if err != nil {
		if err.Error() == "INVALID_ROLE" {
			return nil, status.Error(codes.PermissionDenied, "role cannot access properties")
		}
		return nil, status.Error(codes.Internal, "could not load complexes")
	}
	response := &grpcapi.GetComplexesResponse{Complexes: make([]*grpcapi.Complex, 0, len(complexes))}
	for i := range complexes {
		complex := complexes[i]
		response.Complexes = append(response.Complexes, &complex)
	}
	return response, nil
}

func (server *propertyServer) CreateComplex(_ context.Context, request *grpcapi.CreateComplexRequest) (*grpcapi.CreateComplexResponse, error) {
	if request == nil || request.UserID == 0 || strings.TrimSpace(request.Role) == "" || strings.TrimSpace(request.Name) == "" {
		return nil, status.Error(codes.InvalidArgument, "authenticated user claims and complex name are required")
	}
	if request.Role != "owner" {
		return nil, status.Error(codes.PermissionDenied, "only owners can create complexes")
	}
	complex, err := server.service.CreateComplex(request)
	if err != nil {
		return nil, status.Error(codes.Internal, "could not create complex")
	}
	return &grpcapi.CreateComplexResponse{Complex: complex}, nil
}

func RegisterPropertyServiceServer(registrar grpc.ServiceRegistrar, server PropertyServiceServer) {
	registrar.RegisterService(&PropertyService_ServiceDesc, server)
}

func getComplexesHandler(server any, ctx context.Context, decode func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	request := new(grpcapi.GetComplexesRequest)
	if err := decode(request); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return server.(PropertyServiceServer).GetComplexes(ctx, request)
	}
	info := &grpc.UnaryServerInfo{Server: server, FullMethod: "/" + propertyServiceName + "/GetComplexes"}
	handler := func(ctx context.Context, req any) (any, error) {
		return server.(PropertyServiceServer).GetComplexes(ctx, req.(*grpcapi.GetComplexesRequest))
	}
	return interceptor(ctx, request, info, handler)
}

func createComplexHandler(server any, ctx context.Context, decode func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	request := new(grpcapi.CreateComplexRequest)
	if err := decode(request); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return server.(PropertyServiceServer).CreateComplex(ctx, request)
	}
	info := &grpc.UnaryServerInfo{Server: server, FullMethod: "/" + propertyServiceName + "/CreateComplex"}
	handler := func(ctx context.Context, req any) (any, error) {
		return server.(PropertyServiceServer).CreateComplex(ctx, req.(*grpcapi.CreateComplexRequest))
	}
	return interceptor(ctx, request, info, handler)
}

var PropertyService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: propertyServiceName,
	HandlerType: (*PropertyServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "GetComplexes", Handler: getComplexesHandler},
		{MethodName: "CreateComplex", Handler: createComplexHandler},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "property.proto",
}
