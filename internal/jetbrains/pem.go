package jetbrains

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"time"
)

// ReadPemFile reads a PEM file and returns the decoded DER bytes of its first block.
func ReadPemFile(filepath string) ([]byte, error) {
	certBytes, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("read file %s err %w", filepath, err)
	}
	certBlock, _ := pem.Decode(certBytes)
	return certBlock.Bytes, nil
}

// ReadCertFile reads a PEM-encoded x509 certificate from disk.
func ReadCertFile(filepath string) (*x509.Certificate, error) {
	pemFile, err := ReadPemFile(filepath)
	if err != nil {
		return nil, err
	}
	return x509.ParseCertificate(pemFile)
}

// GenerateRootCertificate signs a child certificate using the given key. The
// name "Root" is historical — the result is actually a leaf certificate signed
// by `issuer`, not a self-signed root.
func GenerateRootCertificate(key *rsa.PrivateKey, subject, issuer string) ([]byte, error) {
	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 80))
	if err != nil {
		return nil, fmt.Errorf("gen serialNumber err %e", err)
	}

	parent := x509.Certificate{
		SerialNumber: serialNumber,
		Subject:      pkix.Name{CommonName: issuer},
		NotBefore:    time.Now().Add(-24 * time.Hour),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}
	template := parent
	template.Subject = pkix.Name{CommonName: subject}

	certBytes, err := x509.CreateCertificate(rand.Reader, &template, &parent, &key.PublicKey, key)
	if err != nil {
		return nil, fmt.Errorf("CreateCertificate err %e", err)
	}
	return certBytes, nil
}
