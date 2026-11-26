package grpcapi

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"os"
	"strings"

	"google.golang.org/grpc/credentials"
)

func LoadServerTLSCredentials() (credentials.TransportCredentials, error) {
	certPath := strings.TrimSpace(os.Getenv("EMPLOYEE_TLS_CERT_FILE"))
	keyPath := strings.TrimSpace(os.Getenv("EMPLOYEE_TLS_KEY_FILE"))
	caPath := strings.TrimSpace(os.Getenv("EMPLOYEE_CLIENT_CA_FILE"))
	if certPath == "" || keyPath == "" || caPath == "" {
		return nil, errors.New("EMPLOYEE_TLS_CERT_FILE, EMPLOYEE_TLS_KEY_FILE, and EMPLOYEE_CLIENT_CA_FILE are required")
	}

	allowedClients := make(map[string]struct{})
	for _, identity := range strings.Split(os.Getenv("EMPLOYEE_ALLOWED_CLIENTS"), ",") {
		identity = strings.ToLower(strings.TrimSpace(identity))
		if identity != "" {
			allowedClients[identity] = struct{}{}
		}
	}
	if len(allowedClients) == 0 {
		return nil, errors.New("EMPLOYEE_ALLOWED_CLIENTS must list at least one authorized client DNS SAN")
	}

	serverCertificate, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, fmt.Errorf("load employee server certificate and key: %w", err)
	}
	caPEM, err := os.ReadFile(caPath)
	if err != nil {
		return nil, fmt.Errorf("read employee client CA certificate: %w", err)
	}
	clientCAs := x509.NewCertPool()
	if !clientCAs.AppendCertsFromPEM(caPEM) {
		return nil, errors.New("employee client CA file contains no valid certificates")
	}
	config := &tls.Config{
		MinVersion: tls.VersionTLS12, Certificates: []tls.Certificate{serverCertificate},
		ClientAuth: tls.RequireAndVerifyClientCert, ClientCAs: clientCAs,
		VerifyConnection: func(state tls.ConnectionState) error {
			if len(state.PeerCertificates) == 0 || len(state.VerifiedChains) == 0 {
				return errors.New("client certificate was not verified")
			}
			for _, dnsName := range state.PeerCertificates[0].DNSNames {
				if _, ok := allowedClients[strings.ToLower(dnsName)]; ok {
					return nil
				}
			}
			return errors.New("client certificate identity is not authorized")
		},
	}
	return credentials.NewTLS(config), nil
}
