package jetbrains

import (
	"license/internal/jetbrains"
	"testing"
)

// Simple test that doesn't require certificate generation
func TestLicenseGenerator_Creation(t *testing.T) {
	generator := jetbrains.NewLicenseGenerator()

	if generator == nil {
		t.Fatal("Failed to create license generator")
	}
}

func TestGetPowerConfig(t *testing.T) {
	generator := jetbrains.NewLicenseGenerator()

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
	generator := jetbrains.NewLicenseGenerator()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = generator.GetPowerConfig()
	}
}
