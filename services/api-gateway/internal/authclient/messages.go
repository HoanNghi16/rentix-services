package authclient

import "fmt"

// These wire-compatible protobuf v1 messages match proto/auth.proto.
type RegisterRequest struct {
	Email    string `protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty"`
	Password string `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
}

func (m *RegisterRequest) Reset()         { *m = RegisterRequest{} }
func (m *RegisterRequest) String() string { return fmt.Sprintf("RegisterRequest{email:%q}", m.Email) }
func (*RegisterRequest) ProtoMessage()    {}

type RegisterResponse struct{}

func (m *RegisterResponse) Reset()         { *m = RegisterResponse{} }
func (m *RegisterResponse) String() string { return "RegisterResponse{}" }
func (*RegisterResponse) ProtoMessage()    {}

type LoginRequest struct {
	Email    string `protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty"`
	Password string `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
}

func (m *LoginRequest) Reset()         { *m = LoginRequest{} }
func (m *LoginRequest) String() string { return fmt.Sprintf("LoginRequest{email:%q}", m.Email) }
func (*LoginRequest) ProtoMessage()    {}

type LoginResponse struct {
	AccessToken  string `protobuf:"bytes,1,opt,name=access_token,json=accessToken,proto3" json:"access_token,omitempty"`
	RefreshToken string `protobuf:"bytes,2,opt,name=refresh_token,json=refreshToken,proto3" json:"refresh_token,omitempty"`
}

func (m *LoginResponse) Reset()         { *m = LoginResponse{} }
func (m *LoginResponse) String() string { return "LoginResponse{tokens:redacted}" }
func (*LoginResponse) ProtoMessage()    {}
