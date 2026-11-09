package grpcapi

import (
	"context"
	"strings"

	"github.com/HoanNghi16/rentix-services/auth-service/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type authServer struct{ service *service.AuthService }

func NewAuthServer(authService *service.AuthService) AuthServiceServer {
	return &authServer{service: authService}
}

func (server *authServer) Register(_ context.Context, request *RegisterRequest) (*RegisterResponse, error) {
	if request == nil || strings.TrimSpace(request.Email) == "" || request.Password == "" {
		return nil, errInvalidAuthRequest
	}
	if err := server.service.Register(request.Email, request.Password); err != nil {
		switch err.Error() {
		case "EMAIL_ALREADY_USED":
			return nil, status.Error(codes.AlreadyExists, "email is already registered")
		case "HASH_FAILED":
			return nil, status.Error(codes.Internal, "could not securely store credentials")
		default:
			return nil, status.Error(codes.Internal, "could not register account")
		}
	}
	return &RegisterResponse{}, nil
}

func (server *authServer) Login(_ context.Context, request *LoginRequest) (*LoginResponse, error) {
	if request == nil || strings.TrimSpace(request.Email) == "" || request.Password == "" {
		return nil, errInvalidAuthRequest
	}
	result, err := server.service.Login(request.Email, request.Password)
	if err != nil {
		switch err.Error() {
		case "USER_NOT_FOUND", "WRONG_PASSWORD":
			return nil, status.Error(codes.Unauthenticated, "invalid email or password")
		default:
			return nil, status.Error(codes.Internal, "could not authenticate account")
		}
	}
	return &LoginResponse{AccessToken: result.AccessToken, RefreshToken: result.RefreshToken}, nil
}
