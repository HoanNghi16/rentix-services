// Command devcert creates local-only certificates for testing auth-service mTLS.
// Never use these self-signed certificates in a deployed environment.
package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"time"
)

func main() {
	const certDir = "certs"
	if err := createDevCertificates(certDir); err != nil {
		fmt.Fprintln(os.Stderr, "create development certificates:", err)
		os.Exit(1)
	}
	fmt.Printf("Created local development certificates in %s\n", certDir)
	fmt.Println("The client certificate identity is api-gateway.rentix.internal")
}

func createDevCertificates(dir string) error {
	if _, err := os.Stat(dir); err == nil {
		return fmt.Errorf("%s already exists; move it aside before generating new certificates", dir)
	} else if !os.IsNotExist(err) {
		return err
	}
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	now := time.Now()
	caKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return err
	}
	caSerial, err := randomSerial()
	if err != nil {
		return err
	}
	caTemplate := &x509.Certificate{
		SerialNumber:          caSerial,
		Subject:               pkix.Name{CommonName: "Rentix local development CA"},
		NotBefore:             now.Add(-5 * time.Minute),
		NotAfter:              now.AddDate(1, 0, 0),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign | x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
		IsCA:                  true,
		MaxPathLenZero:        true,
	}
	caDER, err := x509.CreateCertificate(rand.Reader, caTemplate, caTemplate, &caKey.PublicKey, caKey)
	if err != nil {
		return err
	}
	caCert, err := x509.ParseCertificate(caDER)
	if err != nil {
		return err
	}
	if err := writeCertificate(filepath.Join(dir, "ca.crt"), caDER); err != nil {
		return err
	}
	if err := writePrivateKey(filepath.Join(dir, "ca.key"), caKey); err != nil {
		return err
	}

	serverKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return err
	}
	if err := issueLeaf(dir, "auth-server", "auth-service.rentix.internal", []string{"localhost", "auth-service.rentix.internal"}, []net.IP{net.ParseIP("127.0.0.1")}, x509.ExtKeyUsageServerAuth, caCert, caKey, serverKey, now); err != nil {
		return err
	}

	clientKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return err
	}
	return issueLeaf(dir, "api-gateway", "api-gateway.rentix.internal", []string{"api-gateway.rentix.internal"}, nil, x509.ExtKeyUsageClientAuth, caCert, caKey, clientKey, now)
}

func issueLeaf(dir, filename, commonName string, dnsNames []string, ipAddresses []net.IP, usage x509.ExtKeyUsage, ca *x509.Certificate, caKey *ecdsa.PrivateKey, key *ecdsa.PrivateKey, now time.Time) error {
	serial, err := randomSerial()
	if err != nil {
		return err
	}
	template := &x509.Certificate{
		SerialNumber:          serial,
		Subject:               pkix.Name{CommonName: commonName},
		DNSNames:              dnsNames,
		IPAddresses:           ipAddresses,
		NotBefore:             now.Add(-5 * time.Minute),
		NotAfter:              now.AddDate(0, 3, 0),
		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{usage},
		BasicConstraintsValid: true,
	}
	der, err := x509.CreateCertificate(rand.Reader, template, ca, &key.PublicKey, caKey)
	if err != nil {
		return err
	}
	if err := writeCertificate(filepath.Join(dir, filename+".crt"), der); err != nil {
		return err
	}
	return writePrivateKey(filepath.Join(dir, filename+".key"), key)
}

func randomSerial() (*big.Int, error) {
	return rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
}

func writeCertificate(path string, der []byte) error {
	return writePEM(path, "CERTIFICATE", der)
}

func writePrivateKey(path string, key *ecdsa.PrivateKey) error {
	der, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return err
	}
	return writePEM(path, "PRIVATE KEY", der)
}

func writePEM(path, blockType string, der []byte) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return err
	}
	defer file.Close()
	return pem.Encode(file, &pem.Block{Type: blockType, Bytes: der})
}
