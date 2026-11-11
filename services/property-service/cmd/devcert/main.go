// Command devcert issues a local-only property server certificate from the
// development CA created by auth-service. Do not use it in production.
package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	if err := createCertificate(); err != nil {
		fmt.Fprintln(os.Stderr, "create property development certificate:", err)
		os.Exit(1)
	}
	fmt.Println("Created certs/property-server.crt and certs/property-server.key")
}

func createCertificate() error {
	certPath := envOrDefault("PROPERTY_DEV_CA_CERT_FILE", `..\auth-service\certs\ca.crt`)
	keyPath := envOrDefault("PROPERTY_DEV_CA_KEY_FILE", `..\auth-service\certs\ca.key`)
	caPEM, err := os.ReadFile(certPath)
	if err != nil {
		return fmt.Errorf("read development CA certificate: %w", err)
	}
	caBlock, _ := pem.Decode(caPEM)
	if caBlock == nil || caBlock.Type != "CERTIFICATE" {
		return errors.New("development CA certificate is not valid PEM")
	}
	caCertificate, err := x509.ParseCertificate(caBlock.Bytes)
	if err != nil {
		return err
	}
	keyPEM, err := os.ReadFile(keyPath)
	if err != nil {
		return fmt.Errorf("read development CA private key: %w", err)
	}
	keyBlock, _ := pem.Decode(keyPEM)
	if keyBlock == nil {
		return errors.New("development CA private key is not valid PEM")
	}
	parsedKey, err := x509.ParsePKCS8PrivateKey(keyBlock.Bytes)
	if err != nil {
		return err
	}
	caKey, ok := parsedKey.(*ecdsa.PrivateKey)
	if !ok {
		return errors.New("development CA private key must be ECDSA")
	}

	const certDir = "certs"
	certOutput := filepath.Join(certDir, "property-server.crt")
	keyOutput := filepath.Join(certDir, "property-server.key")
	for _, path := range []string{certOutput, keyOutput} {
		if _, err := os.Stat(path); err == nil {
			return fmt.Errorf("%s already exists; remove it before generating a new certificate", path)
		} else if !os.IsNotExist(err) {
			return err
		}
	}
	if err := os.MkdirAll(certDir, 0700); err != nil {
		return err
	}
	serverKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return err
	}
	serial, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return err
	}
	now := time.Now()
	template := &x509.Certificate{
		SerialNumber:          serial,
		Subject:               pkix.Name{CommonName: "property-service.rentix.internal"},
		DNSNames:              []string{"localhost", "property-service.rentix.internal"},
		IPAddresses:           []net.IP{net.ParseIP("127.0.0.1")},
		NotBefore:             now.Add(-5 * time.Minute),
		NotAfter:              now.AddDate(0, 3, 0),
		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}
	der, err := x509.CreateCertificate(rand.Reader, template, caCertificate, &serverKey.PublicKey, caKey)
	if err != nil {
		return err
	}
	if err := writePEM(certOutput, "CERTIFICATE", der); err != nil {
		return err
	}
	keyDER, err := x509.MarshalPKCS8PrivateKey(serverKey)
	if err != nil {
		return err
	}
	if err := writePEM(keyOutput, "PRIVATE KEY", keyDER); err != nil {
		return err
	}
	return nil
}

func writePEM(path, blockType string, der []byte) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return err
	}
	defer file.Close()
	return pem.Encode(file, &pem.Block{Type: blockType, Bytes: der})
}

func envOrDefault(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}
