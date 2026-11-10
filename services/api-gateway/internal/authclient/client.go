package authclient

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"os"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const authServiceName = "rentix.auth.v1.AuthService"

type Client struct{ connection *grpc.ClientConn }

func NewFromEnv() (*Client, error) {
	target := envOrDefault("AUTH_GRPC_TARGET", "localhost:8081")
	serverName := envOrDefault("AUTH_GRPC_SERVER_NAME", "localhost")
	caPath := envOrDefault("GATEWAY_AUTH_CA_FILE", `..\auth-service\certs\ca.crt`)
	certPath := envOrDefault("GATEWAY_MTLS_CERT_FILE", `..\auth-service\certs\api-gateway.crt`)
	keyPath := envOrDefault("GATEWAY_MTLS_KEY_FILE", `..\auth-service\certs\api-gateway.key`)

	caPEM, err := os.ReadFile(caPath)
	if err != nil {
		return nil, fmt.Errorf("read auth CA certificate: %w", err)
	}
	rootCAs := x509.NewCertPool()
	if !rootCAs.AppendCertsFromPEM(caPEM) {
		return nil, errors.New("auth CA file contains no valid certificates")
	}
	clientCertificate, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, fmt.Errorf("load gateway client certificate and key: %w", err)
	}
	transportCredentials := credentials.NewTLS(&tls.Config{
		MinVersion:   tls.VersionTLS12,
		RootCAs:      rootCAs,
		Certificates: []tls.Certificate{clientCertificate},
		ServerName:   serverName,
	})

	connection, err := grpc.NewClient(target, grpc.WithTransportCredentials(transportCredentials))
	if err != nil {
		return nil, fmt.Errorf("create auth gRPC connection: %w", err)
	}
	return &Client{connection: connection}, nil
}

func (client *Client) Register(ctx context.Context, request *RegisterRequest) (*RegisterResponse, error) {
	response := new(RegisterResponse)
	err := client.connection.Invoke(ctx, "/"+authServiceName+"/Register", request, response)
	return response, err
}

func (client *Client) Login(ctx context.Context, request *LoginRequest) (*LoginResponse, error) {
	response := new(LoginResponse)
	err := client.connection.Invoke(ctx, "/"+authServiceName+"/Login", request, response)
	return response, err
}

func (client *Client) Close() error { return client.connection.Close() }

func envOrDefault(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}
