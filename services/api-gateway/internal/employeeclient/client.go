package employeeclient

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

const employeeServiceName = "rentix.employee.v1.EmployeeService"

type Client struct{ connection *grpc.ClientConn }

func NewFromEnv() (*Client, error) {
	target := envOrDefault("EMPLOYEE_GRPC_TARGET", "localhost:8083")
	serverName := envOrDefault("EMPLOYEE_GRPC_SERVER_NAME", "employee-service.rentix.internal")
	caPath := envOrDefault("GATEWAY_EMPLOYEE_CA_FILE", `..\auth-service\certs\ca.crt`)
	certPath := envOrDefault("GATEWAY_MTLS_CERT_FILE", `..\auth-service\certs\api-gateway.crt`)
	keyPath := envOrDefault("GATEWAY_MTLS_KEY_FILE", `..\auth-service\certs\api-gateway.key`)

	caPEM, err := os.ReadFile(caPath)
	if err != nil {
		return nil, fmt.Errorf("read employee CA certificate: %w", err)
	}
	rootCAs := x509.NewCertPool()
	if !rootCAs.AppendCertsFromPEM(caPEM) {
		return nil, errors.New("employee CA file contains no valid certificates")
	}
	clientCertificate, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, fmt.Errorf("load gateway client certificate and key: %w", err)
	}
	transportCredentials := credentials.NewTLS(&tls.Config{
		MinVersion: tls.VersionTLS12, RootCAs: rootCAs,
		Certificates: []tls.Certificate{clientCertificate}, ServerName: serverName,
	})
	connection, err := grpc.NewClient(target, grpc.WithTransportCredentials(transportCredentials))
	if err != nil {
		return nil, fmt.Errorf("create employee gRPC connection: %w", err)
	}
	return &Client{connection: connection}, nil
}

func (client *Client) GetEmployees(ctx context.Context, request *GetEmployeesRequest) (*GetEmployeesResponse, error) {
	response := new(GetEmployeesResponse)
	err := client.connection.Invoke(ctx, "/"+employeeServiceName+"/GetEmployees", request, response)
	return response, err
}

func (client *Client) Close() error { return client.connection.Close() }

func envOrDefault(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}
