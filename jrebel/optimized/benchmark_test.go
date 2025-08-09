package optimized

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

// BenchmarkOptimizedSign benchmarks the optimized sign function
func BenchmarkOptimizedSign(b *testing.B) {
	controller, err := NewOptimizedLeasesController()
	if err != nil {
		b.Fatalf("Failed to create optimized controller: %v", err)
	}

	clientRandomness := "test-randomness"
	guid := "test-guid-12345"
	offline := true
	validFrom := time.Now().UnixMilli()
	validUntil := validFrom + 180*24*60*60*1000

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		controller.optimizedSign(clientRandomness, guid, offline, validFrom, validUntil)
	}
}

// BenchmarkOptimizedSignWithCache tests performance with cache hits
func BenchmarkOptimizedSignWithCache(b *testing.B) {
	controller, err := NewOptimizedLeasesController()
	if err != nil {
		b.Fatalf("Failed to create optimized controller: %v", err)
	}

	clientRandomness := "test-randomness"
	guid := "test-guid-12345"
	offline := true
	validFrom := time.Now().UnixMilli()
	validUntil := validFrom + 180*24*60*60*1000

	// Pre-populate cache
	controller.optimizedSign(clientRandomness, guid, offline, validFrom, validUntil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		controller.optimizedSign(clientRandomness, guid, offline, validFrom, validUntil)
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

	controller, _ := NewOptimizedLeasesController()
	b.Run("WithPool", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			builder := controller.stringBuilderPool.Get().(*strings.Builder)
			builder.Reset()
			builder.WriteString("test-string")
			_ = builder.String()
			controller.stringBuilderPool.Put(builder)
		}
	})
}

// BenchmarkCacheOperations tests cache performance
func BenchmarkCacheOperations(b *testing.B) {
	controller, err := NewOptimizedLeasesController()
	if err != nil {
		b.Fatalf("Failed to create optimized controller: %v", err)
	}

	keys := make([]string, 1000)
	values := make([]string, 1000)
	for i := 0; i < 1000; i++ {
		keys[i] = "key-" + strconv.Itoa(i)
		values[i] = "value-" + strconv.Itoa(i)
	}

	b.Run("CacheSet", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			key := keys[i%1000]
			value := values[i%1000]
			controller.signatureCache.set(key, value)
		}
	})

	// Pre-populate cache for get operations
	for i := 0; i < 1000; i++ {
		controller.signatureCache.set(keys[i], values[i])
	}

	b.Run("CacheGet", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			key := keys[i%1000]
			controller.signatureCache.get(key)
		}
	})
}

// TestOptimizedSignCorrectness verifies that optimized sign produces the same result as original
func TestOptimizedSignCorrectness(t *testing.T) {
	controller, err := NewOptimizedLeasesController()
	if err != nil {
		t.Fatalf("Failed to create optimized controller: %v", err)
	}

	testCases := []struct {
		clientRandomness string
		guid             string
		offline          bool
		validFrom        int64
		validUntil       int64
	}{
		{"test-randomness", "test-guid", false, 0, 0},
		{"test-randomness", "test-guid", true, 1640995200000, 1656547200000},
		{"another-randomness", "another-guid", true, time.Now().UnixMilli(), time.Now().UnixMilli() + 180*24*60*60*1000},
	}

	for i, tc := range testCases {
		t.Run("TestCase"+strconv.Itoa(i), func(t *testing.T) {
			original := originalSign(tc.clientRandomness, tc.guid, tc.offline, tc.validFrom, tc.validUntil)
			optimized := controller.optimizedSign(tc.clientRandomness, tc.guid, tc.offline, tc.validFrom, tc.validUntil)

			if original != optimized {
				t.Errorf("Signatures don't match.\nOriginal: %s\nOptimized: %s", original, optimized)
			}
		})
	}
}

// TestCacheExpiration tests that cache entries expire correctly
func TestCacheExpiration(t *testing.T) {
	controller, err := NewOptimizedLeasesController()
	if err != nil {
		t.Fatalf("Failed to create optimized controller: %v", err)
	}

	// Set a short TTL for testing
	controller.signatureCache.ttl = 100 * time.Millisecond

	key := "test-key"
	value := "test-signature"
	
	// Set a value
	controller.signatureCache.set(key, value)
	
	// Should be able to get it immediately
	if got := controller.signatureCache.get(key); got != value {
		t.Errorf("Expected %s, got %s", value, got)
	}
	
	// Wait for expiration
	time.Sleep(150 * time.Millisecond)
	
	// Should return empty string after expiration
	if got := controller.signatureCache.get(key); got != "" {
		t.Errorf("Expected empty string after expiration, got %s", got)
	}
}

// BenchmarkParallelAccess tests concurrent access performance
func BenchmarkParallelAccess(b *testing.B) {
	controller, err := NewOptimizedLeasesController()
	if err != nil {
		b.Fatalf("Failed to create optimized controller: %v", err)
	}

	clientRandomness := "test-randomness"
	guid := "test-guid"
	offline := true
	validFrom := time.Now().UnixMilli()
	validUntil := validFrom + 180*24*60*60*1000

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			controller.optimizedSign(clientRandomness, guid, offline, validFrom, validUntil)
		}
	})
}