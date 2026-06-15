package jetbrains

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"license/internal/logger"
	"os"
	"sync"
)

var (
	fake     *FakeCert
	fakeOnce sync.Once
)

// FakeCert is the in-memory representation of the JetBrains certificate chain
// we use to sign licenses. It is a singleton (see GetFake) because the chain
// is process-wide state derived from a single key on disk.
type FakeCert struct {
	CodeRootCert   *x509.Certificate
	ServerRootCert *x509.Certificate
	CodeCert       *x509.Certificate
	ServerCert     *x509.Certificate
	PrivateKey     *rsa.PrivateKey
	PublicKey      *rsa.PublicKey

	ServerUID string

	paths Paths
}

// GetFake returns the singleton FakeCert. Its paths default to the data
// directory; callers that want to override them (e.g. via a CertManager)
// should call SetPaths before any Load/Generate method.
func GetFake() *FakeCert {
	fakeOnce.Do(func() {
		fake = &FakeCert{
			ServerUID: "lemon",
			paths:     DefaultPaths(),
		}
	})
	return fake
}

// SetPaths overrides the certificate paths. Safe to call once during init,
// before any Load/Generate method runs.
func (c *FakeCert) SetPaths(p Paths) {
	c.paths = p
}

// LoadOrGenerate loads the RSA key pair from disk, generating it if missing.
func (c *FakeCert) LoadOrGenerate() error {
	pemFile, err := ReadPemFile(c.paths.PrivateKey)
	if err != nil {
		c.PrivateKey, err = rsa.GenerateKey(rand.Reader, 4096)
		if err != nil {
			return fmt.Errorf("failed to generate private key: %w", err)
		}
		pkcs1PrivateKey := x509.MarshalPKCS1PrivateKey(c.PrivateKey)
		privateKeyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: pkcs1PrivateKey})
		if err = os.WriteFile(c.paths.PrivateKey, privateKeyPEM, 0600); err != nil {
			return fmt.Errorf("failed to write private key: %w", err)
		}
	} else {
		c.PrivateKey, err = x509.ParsePKCS1PrivateKey(pemFile)
		if err != nil {
			return fmt.Errorf("failed to parse private key: %w", err)
		}
	}

	pemFile, err = ReadPemFile(c.paths.PublicKey)
	if err != nil {
		logger.Info("Public key not found, generating new key...")
		pkixPublicKey, err := x509.MarshalPKIXPublicKey(&c.PrivateKey.PublicKey)
		if err != nil {
			return fmt.Errorf("failed to marshal public key: %w", err)
		}
		publicKeyPEM := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pkixPublicKey})
		if err = os.WriteFile(c.paths.PublicKey, publicKeyPEM, 0600); err != nil {
			return fmt.Errorf("failed to write public key: %w", err)
		}
		c.PublicKey = &c.PrivateKey.PublicKey
	} else {
		pub, err := x509.ParsePKIXPublicKey(pemFile)
		if err != nil {
			return fmt.Errorf("failed to parse public key: %w", err)
		}
		var ok bool
		c.PublicKey, ok = pub.(*rsa.PublicKey)
		if !ok {
			return fmt.Errorf("not an RSA public key")
		}
	}

	return nil
}

// LoadRootCert loads the JetBrains code/server CA certificates from disk.
func (c *FakeCert) LoadRootCert() (err error) {
	c.CodeRootCert, err = ReadCertFile(c.paths.CodeRootCert)
	if err != nil {
		return err
	}
	c.ServerRootCert, err = ReadCertFile(c.paths.ServerRootCert)
	return err
}

// LoadCert loads the previously generated leaf code/server certificates.
func (c *FakeCert) LoadCert() (err error) {
	c.CodeCert, err = ReadCertFile(c.paths.CodeCert)
	if err != nil {
		return err
	}
	c.ServerCert, err = ReadCertFile(c.paths.ServerCert)
	return err
}

// GenerateRootCert signs new code/server leaf certificates with the private
// key if they do not already exist on disk.
func (c *FakeCert) GenerateRootCert() (err error) {
	logger.Info("GenerateCodeCert")
	if !fileExists(c.paths.CodeCert) {
		jetCert, err := GenerateRootCertificate(c.PrivateKey, "lemon", c.CodeRootCert.Issuer.CommonName)
		if err != nil {
			return err
		}
		jetCertPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: jetCert})
		if err = os.WriteFile(c.paths.CodeCert, jetCertPEM, 0600); err != nil {
			return err
		}
	}
	logger.Info("GenerateCodeCert done")

	logger.Info("GenerateServerCert")
	if !fileExists(c.paths.ServerCert) {
		subject := fmt.Sprintf("%s.lsrv.jetbrains.com", "lemon")
		lsCert, err := GenerateRootCertificate(c.PrivateKey, subject, c.ServerRootCert.Issuer.CommonName)
		if err != nil {
			return err
		}
		lsCertPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: lsCert})
		if err = os.WriteFile(c.paths.ServerCert, lsCertPEM, 0600); err != nil {
			return err
		}
	}
	logger.Info("GenerateServerCert done")
	return nil
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}
