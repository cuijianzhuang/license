package jetbrains

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha512"
	"encoding/base64"
	"license/internal/logger"
)

// SignWithRsaSha1 signs data with the FakeCert's private key using SHA-1.
func (c *FakeCert) SignWithRsaSha1(data []byte) string {
	hashed := sha1.Sum(data)
	signature, err := rsa.SignPKCS1v15(rand.Reader, c.PrivateKey, crypto.SHA1, hashed[:])
	if err != nil {
		logger.Error("Failed to sign with RSA-SHA1: %v", err)
		return ""
	}
	return base64.StdEncoding.EncodeToString(signature)
}

// SignWithRsaSha512 signs data with the FakeCert's private key using SHA-512.
func (c *FakeCert) SignWithRsaSha512(data []byte) string {
	hashed := sha512.Sum512(data)
	signature, err := rsa.SignPKCS1v15(rand.Reader, c.PrivateKey, crypto.SHA512, hashed[:])
	if err != nil {
		logger.Error("Failed to sign with RSA-SHA512: %v", err)
		return ""
	}
	return base64.StdEncoding.EncodeToString(signature)
}

// CodeCertRawBase64 returns the raw DER bytes of the code certificate, base64-encoded.
func (c *FakeCert) CodeCertRawBase64() string {
	return base64.StdEncoding.EncodeToString(c.CodeCert.Raw)
}

// ServerCertRawBase64 returns the raw DER bytes of the server certificate, base64-encoded.
func (c *FakeCert) ServerCertRawBase64() string {
	return base64.StdEncoding.EncodeToString(c.ServerCert.Raw)
}
