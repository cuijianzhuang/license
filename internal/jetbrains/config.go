package jetbrains

import (
	"time"
)

// DefaultLicenseConfig provides default configuration for license generation
var DefaultLicenseConfig = LicenseConfig{
	DefaultValidDays:   1095, // 3 years
	DefaultGracePeriod: 7,
	DefaultMetadata:    "0120231110PSAA003008",
	DefaultHash:        "51149839/0:-1370131430",
	EnableAI:           true,
	EnableTrial:        false,
}

// productCatalog is the single source of truth for built-in JetBrains product
// codes. Order matters: it determines the order in which fallback products
// appear in a generated license.
var productCatalog = []struct {
	Code string
	Name string
}{
	{"II", "IntelliJ IDEA"},
	{"PS", "PhpStorm"},
	{"AC", "AppCode"},
	{"DB", "DataGrip"},
	{"RM", "RubyMine"},
	{"WS", "WebStorm"},
	{"RD", "Rider"},
	{"CL", "CLion"},
	{"PC", "PyCharm"},
	{"GO", "GoLand"},
	{"DS", "DataSpell"},
	{"DC", "dotCover"},
	{"DPN", "dotPeek"},
	{"DM", "dotMemory"},
}

// DefaultProductCodes returns the built-in JetBrains product codes in catalog order.
func DefaultProductCodes() []string {
	codes := make([]string, len(productCatalog))
	for i, p := range productCatalog {
		codes[i] = p.Code
	}
	return codes
}

// CalculateEffectiveDate calculates the effective date based on configuration
func CalculateEffectiveDate(validDays int) string {
	if validDays <= 0 {
		validDays = DefaultLicenseConfig.DefaultValidDays
	}

	now := time.Now()
	// End of today
	endOfToday := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())
	// Add valid days
	effectiveDate := endOfToday.AddDate(0, 0, validDays)

	return effectiveDate.Format("2006-01-02")
}

// MergeProductCodes merges multiple product code slices and removes duplicates
func MergeProductCodes(codeSets ...[]string) []string {
	seen := make(map[string]bool)
	result := []string{}

	for _, codes := range codeSets {
		for _, code := range codes {
			if !seen[code] {
				seen[code] = true
				result = append(result, code)
			}
		}
	}

	return result
}
