package grpcapi

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const authServiceName = "rentix.auth.v1.AuthService"

type AuthServiceClient interface {
	Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error)
	Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error)
}

type authServiceClient struct{ connection grpc.ClientConnInterface }

func NewAuthServiceClient(connection grpc.ClientConnInterface) AuthServiceClient {
	return &authServiceClient{connection: connection}
}

func (client *authServiceClient) Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error) {
	out := new(RegisterResponse)
	if err := client.connection.Invoke(ctx, "/"+authServiceName+"/Register", in, out, opts...); err != nil {
		return nil, err
	}
	return out, nil
}

func (client *authServiceClient) Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error) {
	out := new(LoginResponse)
	if err := client.connection.Invoke(ctx, "/"+authServiceName+"/Login", in, out, opts...); err != nil {
		return nil, err
	}
	return out, nil
}

type AuthServiceServer interface {
	Register(context.Context, *RegisterRequest) (*RegisterResponse, error)
	Login(context.Context, *LoginRequest) (*LoginResponse, error)
}

func RegisterAuthServiceServer(registrar grpc.ServiceRegistrar, server AuthServiceServer) {
	registrar.RegisterService(&AuthService_ServiceDesc, server)
}

func registerHandler(server any, ctx context.Context, decode func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	request := new(RegisterRequest)
	if err := decode(request); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return server.(AuthServiceServer).Register(ctx, request)
	}
	info := &grpc.UnaryServerInfo{Server: server, FullMethod: "/" + authServiceName + "/Register"}
	handler := func(ctx context.Context, req any) (any, error) {
		return server.(AuthServiceServer).Register(ctx, req.(*RegisterRequest))
	}
	return interceptor(ctx, request, info, handler)
}

func loginHandler(server any, ctx context.Context, decode func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	request := new(LoginRequest)
	if err := decode(request); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return server.(AuthServiceServer).Login(ctx, request)
	}
	info := &grpc.UnaryServerInfo{Server: server, FullMethod: "/" + authServiceName + "/Login"}
	handler := func(ctx context.Context, req any) (any, error) {
		return server.(AuthServiceServer).Login(ctx, req.(*LoginRequest))
	}
	return interceptor(ctx, request, info, handler)
}

var AuthService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: authServiceName,
	HandlerType: (*AuthServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "Register", Handler: registerHandler},
		{MethodName: "Login", Handler: loginHandler},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "auth.proto",
}

var errInvalidAuthRequest = status.Error(codes.InvalidArgument, "email and password are required")
