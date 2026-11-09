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

// LoadServerTLSCredentials requires a trusted client certificate and an
// explicitly authorized DNS SAN before an internal RPC can reach a handler.
func LoadServerTLSCredentials() (credentials.TransportCredentials, error) {
	certPath := strings.TrimSpace(os.Getenv("AUTH_TLS_CERT_FILE"))
	keyPath := strings.TrimSpace(os.Getenv("AUTH_TLS_KEY_FILE"))
	caPath := strings.TrimSpace(os.Getenv("AUTH_CLIENT_CA_FILE"))
	if certPath == "" || keyPath == "" || caPath == "" {
		return nil, errors.New("AUTH_TLS_CERT_FILE, AUTH_TLS_KEY_FILE, and AUTH_CLIENT_CA_FILE are required")
	}

	allowedClients := make(map[string]struct{})
	for _, identity := range strings.Split(os.Getenv("AUTH_ALLOWED_CLIENTS"), ",") {
		identity = strings.ToLower(strings.TrimSpace(identity))
		if identity != "" {
			allowedClients[identity] = struct{}{}
		}
	}
	if len(allowedClients) == 0 {
		return nil, errors.New("AUTH_ALLOWED_CLIENTS must list at least one authorized client DNS SAN")
	}

	serverCertificate, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, fmt.Errorf("load auth server certificate and key: %w", err)
	}
	caPEM, err := os.ReadFile(caPath)
	if err != nil {
		return nil, fmt.Errorf("read client CA certificate: %w", err)
	}
	clientCAs := x509.NewCertPool()
	if ok := clientCAs.AppendCertsFromPEM(caPEM); !ok {
		return nil, errors.New("client CA file contains no valid certificates")
	}

	config := &tls.Config{
		MinVersion:   tls.VersionTLS12,
		Certificates: []tls.Certificate{serverCertificate},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    clientCAs,
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
