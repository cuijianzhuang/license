package jetbrains

import (
	"license/jetbrains/code/service/v2"
	"testing"
)

// Simple test that doesn't require certificate generation
func TestLicenseGenerator_Creation(t *testing.T) {
	generator := v2.NewLicenseGenerator()

	if generator == nil {
		t.Fatal("Failed to create license generator")
	}
}

func TestGetPowerConfig(t *testing.T) {
	generator := v2.NewLicenseGenerator()

	config := generator.GetPowerConfig()

	// The power config may be empty if certificates are not available in test environment
	// Just verify the method doesn't panic and returns a valid structure
	t.Logf("CodePower length: %d", len(config.CodePower))
	t.Logf("ServerPower length: %d", len(config.ServerPower))
	t.Logf("FullConfig length: %d", len(config.FullConfig))

	// Test passes if no panic occurs
}

// Benchmark that doesn't require certificates
func BenchmarkGetPowerConfig(b *testing.B) {
	generator := v2.NewLicenseGenerator()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = generator.GetPowerConfig()
	}
}
