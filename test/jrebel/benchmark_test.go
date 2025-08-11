package jrebel

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"license/jrebel/constant"
	"strconv"
	"strings"
	"testing"
	"time"
)

// Original sign function (copied from the original code for comparison)
func originalSign(clientRandomness, guid string, offline bool, validFrom, validUntil int64) string {
	signatureBase := clientRandomness + ";" + constant.ServerRandomness + ";" + guid + ";" + strconv.FormatBool(offline)
	if offline {
		signatureBase += ";" + strconv.FormatInt(validFrom, 10) + ";" + strconv.FormatInt(validUntil, 10)
	}

	block, _ := pem.Decode([]byte(constant.LeasesPrivateKey))
	if block == nil {
		return ""
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return ""
	}

	hash := sha1.New()
	hash.Write([]byte(signatureBase))
	hashed := hash.Sum(nil)

	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA1, hashed)
	if err != nil {
		return ""
	}

	return base64.StdEncoding.EncodeToString(signature)
}

// BenchmarkOriginalSign benchmarks the original sign function
func BenchmarkOriginalSign(b *testing.B) {
	clientRandomness := "test-randomness"
	guid := "test-guid-12345"
	offline := true
	validFrom := time.Now().UnixMilli()
	validUntil := validFrom + 180*24*60*60*1000

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		originalSign(clientRandomness, guid, offline, validFrom, validUntil)
	}
}

// BenchmarkStringBuilding compares string building methods
func BenchmarkStringBuilding(b *testing.B) {
	clientRandomness := "test-randomness-12345"
	serverRandomness := constant.ServerRandomness
	guid := "test-guid-12345-67890"
	offline := true
	validFrom := int64(1640995200000)
	validUntil := int64(1656547200000)

	b.Run("Original", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			signatureBase := clientRandomness + ";" + serverRandomness + ";" + guid + ";" + strconv.FormatBool(offline)
			if offline {
				signatureBase += ";" + strconv.FormatInt(validFrom, 10) + ";" + strconv.FormatInt(validUntil, 10)
			}
			_ = signatureBase
		}
	})

	b.Run("StringBuilder", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var builder strings.Builder
			builder.WriteString(clientRandomness)
			builder.WriteByte(';')
			builder.WriteString(serverRandomness)
			builder.WriteByte(';')
			builder.WriteString(guid)
			builder.WriteByte(';')
			builder.WriteString(strconv.FormatBool(offline))
			if offline {
				builder.WriteByte(';')
				builder.WriteString(strconv.FormatInt(validFrom, 10))
				builder.WriteByte(';')
				builder.WriteString(strconv.FormatInt(validUntil, 10))
			}
			_ = builder.String()
		}
	})

	b.Run("BytesBuffer", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var buf bytes.Buffer
			buf.WriteString(clientRandomness)
			buf.WriteByte(';')
			buf.WriteString(serverRandomness)
			buf.WriteByte(';')
			buf.WriteString(guid)
			buf.WriteByte(';')
			buf.WriteString(strconv.FormatBool(offline))
			if offline {
				buf.WriteByte(';')
				buf.WriteString(strconv.FormatInt(validFrom, 10))
				buf.WriteByte(';')
				buf.WriteString(strconv.FormatInt(validUntil, 10))
			}
			_ = buf.String()
		}
	})
}

// BenchmarkMemoryAllocation tests memory allocation patterns
func BenchmarkMemoryAllocation(b *testing.B) {
	b.Run("WithoutPool", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			builder := &strings.Builder{}
			builder.WriteString("test-string")
			_ = builder.String()
		}
	})

	b.Run("WithPool", func(b *testing.B) {
		pool := &strings.Builder{}
		for i := 0; i < b.N; i++ {
			pool.Reset()
			pool.WriteString("test-string")
			_ = pool.String()
		}
	})
}