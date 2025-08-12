package v2

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"license/config"
	jetbrainsConfig "license/jetbrains/config"
	"license/jetbrains/types"
	"license/jetbrains/util"
	"license/logger"
	"math/big"
	"sync"
	"time"
)

// LicenseGenerator handles license generation
type LicenseGenerator struct {
	fakeCert *util.FakeCert
	mu       sync.RWMutex
	cache    map[string]*cachedLicense
}

// cachedLicense represents a cached license
type cachedLicense struct {
	license   string
	timestamp time.Time
}

// NewLicenseGenerator creates a new license generator
func NewLicenseGenerator() *LicenseGenerator {
	return &LicenseGenerator{
		fakeCert: util.GetFake(),
		cache:    make(map[string]*cachedLicense),
	}
}

// GenerateLicense generates a JetBrains license
func (g *LicenseGenerator) GenerateLicense(req types.GenerateLicenseRequest) (*types.GenerateLicenseResponse, error) {
	// Validate request
	if req.LicenseeName == "" {
		return nil, fmt.Errorf("license name is required")
	}

	// Calculate effective date
	effectiveDate := req.EffectiveDate
	if effectiveDate == "" {
		effectiveDate = jetbrainsConfig.CalculateEffectiveDate(req.ValidDays)
	}

	// Get product codes
	codes := g.getProductCodes(req.Codes)
	if len(codes) == 0 {
		return nil, fmt.Errorf("no product codes available")
	}

	// Generate license ID
	licenseID, err := g.generateLicenseID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate license ID: %w", err)
	}

	// Check cache
	cacheKey := fmt.Sprintf("%s:%s:%v", req.LicenseeName, effectiveDate, codes)
	if cached := g.getFromCache(cacheKey); cached != "" {
		logger.Info("Using cached license for: " + req.LicenseeName)
		return g.buildResponse(cached, licenseID, effectiveDate), nil
	}

	// Build products list
	products := g.buildProducts(codes, effectiveDate)

	// Create license part
	licensePart := types.LicensePart{
		LicenseID:         licenseID,
		LicenseeName:      req.LicenseeName,
		Products:          products,
		AssigneeName:      "",
		Metadata:          jetbrainsConfig.DefaultLicenseConfig.DefaultMetadata,
		Hash:              jetbrainsConfig.DefaultLicenseConfig.DefaultHash,
		GracePeriodDays:   jetbrainsConfig.DefaultLicenseConfig.DefaultGracePeriod,
		AutoProlongated:   true,
		IsAutoProlongated: true,
		Trial:             jetbrainsConfig.DefaultLicenseConfig.EnableTrial,
		AiAllowed:         jetbrainsConfig.DefaultLicenseConfig.EnableAI,
	}

	// Generate activation code
	activationCode, err := g.generateActivationCode(licensePart)
	if err != nil {
		return nil, fmt.Errorf("failed to generate activation code: %w", err)
	}

	// Cache the result
	g.saveToCache(cacheKey, activationCode)

	// Build response
	return g.buildResponse(activationCode, licenseID, effectiveDate), nil
}

// getProductCodes retrieves and merges product codes
func (g *LicenseGenerator) getProductCodes(requestedCodes []string) []string {
	var allCodes []string

	// Add requested codes
	if len(requestedCodes) > 0 {
		allCodes = append(allCodes, requestedCodes...)
	}

	// Add codes from database
	dbCodes := g.getCodesFromDatabase()
	allCodes = append(allCodes, dbCodes...)

	// Add default product codes
	allCodes = append(allCodes, jetbrainsConfig.DefaultProductCodes...)

	// Remove duplicates
	return jetbrainsConfig.MergeProductCodes(allCodes)
}

// getCodesFromDatabase retrieves codes from database
func (g *LicenseGenerator) getCodesFromDatabase() []string {
	var codes []string

	// Skip database access if DB is not initialized
	if config.DB == nil {
		logger.Info("Database not initialized, skipping database codes")
		return codes
	}

	// Get product codes
	productService := NewProductService()
	products, err := productService.GetAll()
	if err == nil {
		for _, product := range products {
			codes = append(codes, product.ProductCode)
		}
	} else {
		logger.Info("Failed to get products from database: " + err.Error())
	}

	// Get plugin codes
	pluginService := NewPluginService()
	plugins, err := pluginService.GetAll()
	if err == nil {
		for _, plugin := range plugins {
			codes = append(codes, plugin.PluginCode)
		}
	} else {
		logger.Info("Failed to get plugins from database: " + err.Error())
	}

	return codes
}

// buildProducts creates product entries
func (g *LicenseGenerator) buildProducts(codes []string, effectiveDate string) []types.Product {
	var products []types.Product

	for _, code := range codes {
		products = append(products, types.Product{
			Code:         code,
			FallbackDate: effectiveDate,
			PaidUpTo:     effectiveDate,
			Extended:     true,
		})
	}

	return products
}

// generateActivationCode creates the final activation code
func (g *LicenseGenerator) generateActivationCode(licensePart types.LicensePart) (string, error) {
	// Marshal license part to JSON
	licenseJSON, err := json.Marshal(licensePart)
	if err != nil {
		return "", fmt.Errorf("failed to marshal license: %w", err)
	}

	// Encode to base64
	licenseBase64 := base64.StdEncoding.EncodeToString(licenseJSON)

	// Sign the license
	signature, err := g.signLicense(licenseJSON)
	if err != nil {
		return "", fmt.Errorf("failed to sign license: %w", err)
	}

	// Get certificate
	cert := g.fakeCert.CodeCert
	if cert == nil {
		return "", fmt.Errorf("certificate not available")
	}
	certBase64 := base64.StdEncoding.EncodeToString(cert.Raw)

	// Combine parts
	activationCode := fmt.Sprintf("%s-%s-%s-%s",
		licensePart.LicenseID,
		licenseBase64,
		signature,
		certBase64,
	)

	return activationCode, nil
}

// signLicense signs the license data
func (g *LicenseGenerator) signLicense(data []byte) (string, error) {
	if g.fakeCert.PrivateKey == nil {
		return "", fmt.Errorf("private key not available")
	}

	hashed := sha1.Sum(data)
	signature, err := rsa.SignPKCS1v15(rand.Reader, g.fakeCert.PrivateKey, crypto.SHA1, hashed[:])
	if err != nil {
		return "", fmt.Errorf("failed to sign data: %w", err)
	}

	return base64.StdEncoding.EncodeToString(signature), nil
}

// generateLicenseID generates a random license ID
func (g *LicenseGenerator) generateLicenseID() (string, error) {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 10

	result := make([]byte, length)
	for i := range result {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result[i] = charset[n.Int64()]
	}

	return string(result), nil
}

// buildResponse creates the API response
func (g *LicenseGenerator) buildResponse(activationCode, licenseID, expiresAt string) *types.GenerateLicenseResponse {
	// Generate power config
	powerConfig := g.generatePowerConfig()

	return &types.GenerateLicenseResponse{
		ActivationCode: activationCode,
		PowerConfig:    powerConfig,
		LicenseID:      licenseID,
		ExpiresAt:      expiresAt,
		GeneratedAt:    time.Now(),
	}
}

// generatePowerConfig generates the power.conf content
func (g *LicenseGenerator) generatePowerConfig() string {
	if g.fakeCert.CodeCert == nil || g.fakeCert.CodeRootCert == nil {
		return ""
	}

	codePower := util.GeneratePowerResult(g.fakeCert.CodeCert, g.fakeCert.CodeRootCert)

	return fmt.Sprintf("[Result]\n; Lemon active by code\n%s", codePower)
}

// GetPowerConfig returns the full power configuration
func (g *LicenseGenerator) GetPowerConfig() types.PowerConfigResponse {
	var codePower, serverPower string

	if g.fakeCert.CodeCert != nil && g.fakeCert.CodeRootCert != nil {
		codePower = util.GeneratePowerResult(g.fakeCert.CodeCert, g.fakeCert.CodeRootCert)
	}

	if g.fakeCert.ServerCert != nil && g.fakeCert.ServerRootCert != nil {
		serverPower = util.GeneratePowerResult(g.fakeCert.ServerCert, g.fakeCert.ServerRootCert)
	}

	fullConfig := fmt.Sprintf("[Result]\n; Lemon active by code\n%s\n[Result]\n; Lemon active by server\n%s",
		codePower, serverPower)

	return types.PowerConfigResponse{
		CodePower:   codePower,
		ServerPower: serverPower,
		FullConfig:  fullConfig,
	}
}

// Cache management methods

func (g *LicenseGenerator) getFromCache(key string) string {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if entry, ok := g.cache[key]; ok {
		// Cache expires after 1 hour
		if time.Since(entry.timestamp) < time.Hour {
			return entry.license
		}
	}
	return ""
}

func (g *LicenseGenerator) saveToCache(key, license string) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.cache[key] = &cachedLicense{
		license:   license,
		timestamp: time.Now(),
	}

	// Clean old entries if cache is too large
	if len(g.cache) > 100 {
		g.cleanCache()
	}
}

func (g *LicenseGenerator) cleanCache() {
	now := time.Now()
	for key, entry := range g.cache {
		if now.Sub(entry.timestamp) > time.Hour {
			delete(g.cache, key)
		}
	}
}
