package propertyclient

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

const propertyServiceName = "rentix.property.v1.PropertyService"

type Client struct{ connection *grpc.ClientConn }

func NewFromEnv() (*Client, error) {
	target := envOrDefault("PROPERTY_GRPC_TARGET", "localhost:8082")
	serverName := envOrDefault("PROPERTY_GRPC_SERVER_NAME", "localhost")
	caPath := envOrDefault("GATEWAY_PROPERTY_CA_FILE", `..\auth-service\certs\ca.crt`)
	certPath := envOrDefault("GATEWAY_MTLS_CERT_FILE", `..\auth-service\certs\api-gateway.crt`)
	keyPath := envOrDefault("GATEWAY_MTLS_KEY_FILE", `..\auth-service\certs\api-gateway.key`)

	caPEM, err := os.ReadFile(caPath)
	if err != nil {
		return nil, fmt.Errorf("read property CA certificate: %w", err)
	}
	rootCAs := x509.NewCertPool()
	if !rootCAs.AppendCertsFromPEM(caPEM) {
		return nil, errors.New("property CA file contains no valid certificates")
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
		return nil, fmt.Errorf("create property gRPC connection: %w", err)
	}
	return &Client{connection: connection}, nil
}

func (client *Client) GetComplexes(ctx context.Context, request *GetComplexesRequest) (*GetComplexesResponse, error) {
	response := new(GetComplexesResponse)
	err := client.connection.Invoke(ctx, "/"+propertyServiceName+"/GetComplexes", request, response)
	return response, err
}

func (client *Client) GetRooms(ctx context.Context, request *GetRoomsRequest) (*GetRoomsResponse, error) {
	response := new(GetRoomsResponse)
	err := client.connection.Invoke(ctx, "/"+propertyServiceName+"/GetRooms", request, response)
	return response, err
}

func (client *Client) CreateComplex(ctx context.Context, request *CreateComplexRequest) (*CreateComplexResponse, error) {
	response := new(CreateComplexResponse)
	err := client.connection.Invoke(ctx, "/"+propertyServiceName+"/CreateComplex", request, response)
	return response, err
}

func (client *Client) Close() error { return client.connection.Close() }

func envOrDefault(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}
