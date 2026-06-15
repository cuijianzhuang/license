package types

import (
	"strings"
	"time"
)

// LicensePart represents the license information structure
type LicensePart struct {
	LicenseID         string    `json:"licenseId"`
	LicenseeName      string    `json:"licenseeName"`
	Products          []Product `json:"products"`
	AssigneeName      string    `json:"assigneeName"`
	Metadata          string    `json:"metadata"`
	Hash              string    `json:"hash"`
	GracePeriodDays   int       `json:"gracePeriodDays"`
	AutoProlongated   bool      `json:"autoProlongated"`
	IsAutoProlongated bool      `json:"isAutoProlongated"`
	Trial             bool      `json:"trial"`
	AiAllowed         bool      `json:"aiAllowed"`
}

// Product represents a JetBrains product license entry
type Product struct {
	Code         string `json:"code"`
	FallbackDate string `json:"fallbackDate"`
	PaidUpTo     string `json:"paidUpTo"`
	Extended     bool   `json:"extended"`
}

// GenerateLicenseRequest represents the API request for license generation
type GenerateLicenseRequest struct {
	LicenseeName  string `json:"licenseeName" form:"licenseeName" binding:"required,min=1"`
	EffectiveDate string `json:"effectiveDate" form:"effectiveDate"`
	Codes         string `json:"codes" form:"codes"`
	ValidDays     int    `json:"validDays" form:"validDays"`
}

// ParseCodes splits Codes by comma, trims whitespace and deduplicates
func (r *GenerateLicenseRequest) ParseCodes() []string {
	if r.Codes == "" {
		return nil
	}
	set := make(map[string]struct{})
	var result []string
	for _, p := range strings.Split(r.Codes, ",") {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if _, exists := set[p]; !exists {
			set[p] = struct{}{}
			result = append(result, p)
		}
	}
	return result
}

// GenerateLicenseResponse represents the API response for license generation
type GenerateLicenseResponse struct {
	ActivationCode string    `json:"activationCode"`
	PowerConfig    string    `json:"powerConfig"`
	LicenseID      string    `json:"licenseId"`
	ExpiresAt      string    `json:"expiresAt"`
	GeneratedAt    time.Time `json:"generatedAt"`
}

// PowerConfigResponse represents the power.conf response
type PowerConfigResponse struct {
	CodePower   string `json:"codePower"`
	ServerPower string `json:"serverPower"`
	FullConfig  string `json:"fullConfig"`
}

// PluginInfo represents plugin information
type PluginInfo struct {
	ID       uint64 `json:"id"`
	PluginID uint64 `json:"pluginId"`
	Name     string `json:"name"`
	Code     string `json:"code"`
	Detail   string `json:"detail,omitempty"`
}

// ProductInfo represents product information
type ProductInfo struct {
	ID     uint64 `json:"id"`
	Name   string `json:"name"`
	Code   string `json:"code"`
	Detail string `json:"detail,omitempty"`
}

// LicenseConfig represents license generation configuration
type LicenseConfig struct {
	DefaultValidDays   int
	DefaultGracePeriod int
	DefaultMetadata    string
	DefaultHash        string
	EnableAI           bool
	EnableTrial        bool
}
