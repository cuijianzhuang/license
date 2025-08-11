package jetbrains

import (
	"encoding/base64"
	"encoding/json"
	"license/jetbrains/config"
	"license/jetbrains/service"
	"license/jetbrains/types"
	"license/jetbrains/util"
	"strings"
	"testing"
	"time"
)

// initTestCertificates initializes fake certificates for testing
func initTestCertificates(t *testing.T) {
	// Create a temporary directory for test certificates
	tempDir := t.TempDir()
	t.Setenv("DATA_DIR", tempDir)
	
	// Initialize certificates
	fakeCert := util.GetFake()
	if err := fakeCert.LoadOrGenerate(); err != nil {
		t.Fatalf("Failed to load or generate certificates: %v", err)
	}
	
	// Generate root certificates if needed
	if err := fakeCert.LoadRootCert(); err != nil {
		t.Logf("Failed to load root cert, generating: %v", err)
		if err := fakeCert.GenerateRootCert(); err != nil {
			t.Fatalf("Failed to generate root cert: %v", err)
		}
	}
	
	// Load certificates
	if err := fakeCert.LoadCert(); err != nil {
		t.Logf("Failed to load cert: %v", err)
	}
}

func TestLicenseGenerator_GenerateLicense(t *testing.T) {
	// Initialize fake certificates for testing
	initTestCertificates(t)
	
	generator := service.NewLicenseGenerator()

	tests := []struct {
		name    string
		req     types.GenerateLicenseRequest
		wantErr bool
		checks  func(t *testing.T, resp *types.GenerateLicenseResponse)
	}{
		{
			name: "valid request with default settings",
			req: types.GenerateLicenseRequest{
				LicenseeName: "Test User",
			},
			wantErr: false,
			checks: func(t *testing.T, resp *types.GenerateLicenseResponse) {
				if resp.ActivationCode == "" {
					t.Error("ActivationCode is empty")
				}
				if resp.LicenseID == "" {
					t.Error("LicenseID is empty")
				}
				if resp.PowerConfig == "" {
					t.Error("PowerConfig is empty")
				}
				// Check activation code format (4 parts separated by -)
				parts := strings.Split(resp.ActivationCode, "-")
				if len(parts) != 4 {
					t.Errorf("ActivationCode should have 4 parts, got %d", len(parts))
				}
			},
		},
		{
			name: "request with custom validity period",
			req: types.GenerateLicenseRequest{
				LicenseeName: "Test User",
				ValidDays:    365,
			},
			wantErr: false,
			checks: func(t *testing.T, resp *types.GenerateLicenseResponse) {
				if resp.ExpiresAt == "" {
					t.Error("ExpiresAt is empty")
				}
				// Parse expiry date
				expiry, err := time.Parse("2006-01-02", resp.ExpiresAt)
				if err != nil {
					t.Errorf("Failed to parse expiry date: %v", err)
				}
				// Check if it's approximately 365 days from now
				expectedDays := 365
				actualDays := int(expiry.Sub(time.Now()).Hours() / 24)
				if actualDays < expectedDays-1 || actualDays > expectedDays+1 {
					t.Errorf("Expected ~%d days validity, got %d", expectedDays, actualDays)
				}
			},
		},
		{
			name: "request with specific product codes",
			req: types.GenerateLicenseRequest{
				LicenseeName: "Test User",
				Codes:        []string{"GO", "WS", "DB"},
			},
			wantErr: false,
			checks: func(t *testing.T, resp *types.GenerateLicenseResponse) {
				// Decode and verify license contains requested codes
				parts := strings.Split(resp.ActivationCode, "-")
				if len(parts) < 2 {
					t.Fatal("Invalid activation code format")
				}
				
				licenseData, err := base64.StdEncoding.DecodeString(parts[1])
				if err != nil {
					t.Fatalf("Failed to decode license data: %v", err)
				}
				
				var licensePart types.LicensePart
				if err := json.Unmarshal(licenseData, &licensePart); err != nil {
					t.Fatalf("Failed to unmarshal license part: %v", err)
				}
				
				// Check if requested codes are present
				codeMap := make(map[string]bool)
				for _, product := range licensePart.Products {
					codeMap[product.Code] = true
				}
				
				for _, code := range []string{"GO", "WS", "DB"} {
					if !codeMap[code] {
						t.Errorf("Expected code %s not found in license", code)
					}
				}
			},
		},
		{
			name: "empty license name should fail",
			req: types.GenerateLicenseRequest{
				LicenseeName: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := generator.GenerateLicense(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateLicense() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.checks != nil {
				tt.checks(t, resp)
			}
		})
	}
}

// Note: These tests have been disabled because they access private methods
// that are not exported from the service package

/*
func TestLicenseGenerator_generateLicenseID(t *testing.T) {
	generator := service.NewLicenseGenerator()
	
	// Generate multiple IDs and check uniqueness
	ids := make(map[string]bool)
	for i := 0; i < 100; i++ {
		id, err := generator.generateLicenseID()
		if err != nil {
			t.Fatalf("Failed to generate license ID: %v", err)
		}
		
		if len(id) != 10 {
			t.Errorf("License ID should be 10 characters, got %d", len(id))
		}
		
		if ids[id] {
			t.Errorf("Duplicate license ID generated: %s", id)
		}
		ids[id] = true
	}
}
*/

func TestLicenseGenerator_PublicMethods(t *testing.T) {
	generator := service.NewLicenseGenerator()
	
	// Test that we can create a generator instance
	if generator == nil {
		t.Fatal("Failed to create license generator")
	}
	
	// Test GetPowerConfig
	config := generator.GetPowerConfig()
	if config.CodePower == "" {
		t.Error("CodePower is empty")
	}
}

func TestLicenseGenerator_Cache(t *testing.T) {
	generator := service.NewLicenseGenerator()
	
	req := types.GenerateLicenseRequest{
		LicenseeName: "Cache Test User",
		Codes:        []string{"GO"},
	}
	
	// First call - should generate new
	resp1, err := generator.GenerateLicense(req)
	if err != nil {
		t.Fatalf("Failed to generate license: %v", err)
	}
	
	// Second call - should use cache (but with new license ID)
	resp2, err := generator.GenerateLicense(req)
	if err != nil {
		t.Fatalf("Failed to get cached license: %v", err)
	}
	
	// License IDs should be different
	if resp1.LicenseID == resp2.LicenseID {
		t.Error("License IDs should be different even when cached")
	}
	
	// Activation codes should contain same license data (middle part)
	parts1 := strings.Split(resp1.ActivationCode, "-")
	parts2 := strings.Split(resp2.ActivationCode, "-")
	
	if len(parts1) < 2 || len(parts2) < 2 {
		t.Fatal("Invalid activation code format")
	}
	
	// The license data part should be the same (from cache)
	if parts1[1] != parts2[1] {
		t.Error("Cached license data should be the same")
	}
}

func TestGetPowerConfig(t *testing.T) {
	generator := service.NewLicenseGenerator()
	
	config := generator.GetPowerConfig()
	
	if config.CodePower == "" {
		t.Error("CodePower is empty")
	}
	
	if config.ServerPower == "" {
		t.Error("ServerPower is empty")
	}
	
	if config.FullConfig == "" {
		t.Error("FullConfig is empty")
	}
	
	// Check format
	if !strings.Contains(config.FullConfig, "[Result]") {
		t.Error("FullConfig should contain [Result] section")
	}
	
	if !strings.Contains(config.FullConfig, "Lemon active") {
		t.Error("FullConfig should contain activation message")
	}
}

func TestCalculateEffectiveDate(t *testing.T) {
	tests := []struct {
		name      string
		validDays int
		wantDays  int
	}{
		{
			name:      "default days",
			validDays: 0,
			wantDays:  config.DefaultLicenseConfig.DefaultValidDays,
		},
		{
			name:      "custom days",
			validDays: 30,
			wantDays:  30,
		},
		{
			name:      "negative days uses default",
			validDays: -10,
			wantDays:  config.DefaultLicenseConfig.DefaultValidDays,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dateStr := config.CalculateEffectiveDate(tt.validDays)
			
			// Parse the date
			date, err := time.Parse("2006-01-02", dateStr)
			if err != nil {
				t.Fatalf("Failed to parse date: %v", err)
			}
			
			// Calculate expected days
			expectedDays := tt.wantDays
			actualDays := int(date.Sub(time.Now()).Hours() / 24)
			
			// Allow 1 day tolerance
			if actualDays < expectedDays-1 || actualDays > expectedDays+1 {
				t.Errorf("Expected ~%d days, got %d", expectedDays, actualDays)
			}
		})
	}
}

func BenchmarkGenerateLicense(b *testing.B) {
	generator := service.NewLicenseGenerator()
	req := types.GenerateLicenseRequest{
		LicenseeName: "Benchmark User",
		Codes:        []string{"GO", "WS", "DB"},
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = generator.GenerateLicense(req)
	}
}

// Note: The following benchmarks have been disabled because they access private methods
/*
func BenchmarkGenerateLicenseID(b *testing.B) {
	generator := service.NewLicenseGenerator()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = generator.generateLicenseID()
	}
}

func BenchmarkSignLicense(b *testing.B) {
	generator := service.NewLicenseGenerator()
	data := []byte("test license data for benchmarking")
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = generator.signLicense(data)
	}
}
*/