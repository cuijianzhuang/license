package v2

import (
	"fmt"
	v2 "license/jetbrains/code/service/v2"
	"net/http"
	"time"

	"license/jetbrains/types"
	"license/logger"
	v1 "license/v1"

	"github.com/gin-gonic/gin"
)

// Controller handles JetBrains license API endpoints
type Controller struct {
	generator      *v2.LicenseGenerator
	productService *v2.ProductService
	pluginService  *v2.PluginService
}

// NewController creates a new JetBrains controller
func NewController() *Controller {
	return &Controller{
		generator:      v2.NewLicenseGenerator(),
		productService: v2.NewProductService(),
		pluginService:  v2.NewPluginService(),
	}
}

// GenerateLicense handles license generation requests
func (c *Controller) GenerateLicense(ctx *gin.Context) {
	var req types.GenerateLicenseRequest

	// Try to bind from JSON body
	if err := ctx.ShouldBind(&req); err != nil {
		logger.Error("ShouldBind failed, fallback to query parameters: ", err)
	}

	// Fallback to query parameters for empty fields
	if req.LicenseeName == "" {
		req.LicenseeName = ctx.Query("licenseeName")
	}
	if req.EffectiveDate == "" {
		req.EffectiveDate = ctx.Query("effectiveDate")
	}
	if req.Codes == "" {
		req.Codes = ctx.Query("codes")
	}
	if req.ValidDays == 0 {
		if validDaysStr := ctx.Query("validDays"); validDaysStr != "" {
			var validDays int
			if _, err := fmt.Sscanf(validDaysStr, "%d", &validDays); err == nil {
				req.ValidDays = validDays
			}
		}
	}

	// Validate required fields
	if req.LicenseeName == "" {
		v1.HandleError(ctx, http.StatusBadRequest, "License name is required")
		return
	}

	// Log request
	logger.Info("Generating JetBrains license for: " + req.LicenseeName)

	// Generate license
	response, err := c.generator.GenerateLicense(req)
	if err != nil {
		logger.Error("Failed to generate license", err)
		v1.HandleError(ctx, http.StatusInternalServerError, "Failed to generate license")
		return
	}

	v1.HandleSuccess(ctx, response)
}

// GetPowerConfig returns the power.conf configuration
func (c *Controller) GetPowerConfig(ctx *gin.Context) {
	powerConfig := c.generator.GetPowerConfig()

	// Check output format
	format := ctx.Query("format")
	if format == "text" || format == "raw" {
		ctx.String(http.StatusOK, powerConfig.FullConfig)
		return
	}

	v1.HandleSuccess(ctx, powerConfig)
}

// FetchProductsLatest fetches the latest products
func (c *Controller) FetchProductsLatest(ctx *gin.Context) {
	// Run in background
	go func() {
		if err := c.productService.FetchLatest(); err != nil {
			logger.Error("Failed to fetch latest products", err)
		}
	}()

	v1.HandleSuccess(ctx, gin.H{
		"message": "Fetching latest products in background",
		"status":  "processing",
	})
}

// FetchPluginsLatest fetches the latest plugins
func (c *Controller) FetchPluginsLatest(ctx *gin.Context) {
	// Run in background
	go func() {
		if err := c.pluginService.FetchLatest(); err != nil {
			logger.Error("Failed to fetch latest plugins", err)
		}
	}()

	v1.HandleSuccess(ctx, gin.H{
		"message": "Fetching latest plugins in background",
		"status":  "processing",
	})
}

// GetProducts returns all available products
func (c *Controller) GetProducts(ctx *gin.Context) {
	products, err := c.productService.GetAll()
	if err != nil {
		logger.Error("Failed to get products", err)
		v1.HandleError(ctx, http.StatusInternalServerError, "Failed to get products")
		return
	}

	// Convert to response format
	var response []types.ProductInfo
	for _, product := range products {
		response = append(response, types.ProductInfo{
			ID:     product.ID,
			Name:   product.ProductName,
			Code:   product.ProductCode,
			Detail: product.ProductDetail,
		})
	}

	v1.HandleSuccess(ctx, response)
}

// GetPlugins returns all available plugins
func (c *Controller) GetPlugins(ctx *gin.Context) {
	plugins, err := c.pluginService.GetAll()
	if err != nil {
		logger.Error("Failed to get plugins", err)
		v1.HandleError(ctx, http.StatusInternalServerError, "Failed to get plugins")
		return
	}

	// Convert to response format
	var response []types.PluginInfo
	for _, plugin := range plugins {
		response = append(response, types.PluginInfo{
			ID:       plugin.ID,
			PluginID: plugin.PluginID,
			Name:     plugin.PluginName,
			Code:     plugin.PluginCode,
			Detail:   plugin.PluginApiDetail,
		})
	}

	v1.HandleSuccess(ctx, response)
}

// HealthCheck provides a health check endpoint
func (c *Controller) HealthCheck(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "jetbrains",
		"time":    time.Now().Unix(),
	})
}

// ServerController handles JetBrains server API endpoints
type ServerController struct {
	generator *v2.LicenseGenerator
}

// NewServerController creates a new server controller
func NewServerController() *ServerController {
	return &ServerController{
		generator: v2.NewLicenseGenerator(),
	}
}

// LicenseServerRule returns the license server rules
func (sc *ServerController) LicenseServerRule(ctx *gin.Context) {
	powerConfig := sc.generator.GetPowerConfig()
	ctx.String(http.StatusOK, powerConfig.FullConfig)
}
