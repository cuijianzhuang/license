package types

import (
	"strings"
	"time"
)

// LicensePart represents the license information structure
type LicensePart struct {
	LicenseID         string            `json:"licenseId"`
	LicenseeName      string            `json:"licenseeName"`
	Products          []LicensedProduct `json:"products"`
	AssigneeName      string            `json:"assigneeName"`
	Metadata          string            `json:"metadata"`
	Hash              string            `json:"hash"`
	GracePeriodDays   int               `json:"gracePeriodDays"`
	AutoProlongated   bool              `json:"autoProlongated"`
	IsAutoProlongated bool              `json:"isAutoProlongated"`
	Trial             bool              `json:"trial"`
	AiAllowed         bool              `json:"aiAllowed"`
}

// LicensedProduct is a single product entry inside a generated license.
// (Distinct from Product below, which models the products catalog row.)
type LicensedProduct struct {
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

// Plugin is the GORM-mapped JetBrains plugin record. The json tags also drive
// the GET /jetbrains/plugins response shape.
type Plugin struct {
	ID              uint64 `gorm:"primaryKey;column:id" json:"id"`
	PluginID        uint64 `gorm:"column:plugin_id" json:"pluginId"`
	PluginName      string `gorm:"column:plugin_name" json:"name"`
	PluginCode      string `gorm:"column:plugin_code;uniqueIndex:idx_plugin_code" json:"code"`
	PluginApiDetail string `gorm:"column:plugin_api_detail" json:"detail,omitempty"`
}

// TableName overrides the GORM default so plugins live in sys_jetbrains_paid_plugin.
func (Plugin) TableName() string {
	return "sys_jetbrains_paid_plugin"
}

// Product is the GORM-mapped JetBrains product record. The json tags also drive
// the GET /jetbrains/products response shape.
type Product struct {
	ID            uint64 `gorm:"primaryKey;column:id" json:"id"`
	ProductName   string `gorm:"column:product_name" json:"name"`
	ProductCode   string `gorm:"column:product_code;uniqueIndex:idx_product_code" json:"code"`
	ProductDetail string `gorm:"column:product_detail" json:"detail,omitempty"`
}

// TableName overrides the GORM default so products live in sys_jetbrains_product.
func (Product) TableName() string {
	return "sys_jetbrains_product"
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
