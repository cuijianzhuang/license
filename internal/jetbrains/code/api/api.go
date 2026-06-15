package api

import (
	"license/internal/jetbrains/code/service"
	"net/http"
	"time"

	"license/internal/jetbrains/types"
	"license/internal/logger"
	"license/internal/httpx"

	"github.com/gin-gonic/gin"
)

// Controller handles JetBrains license API endpoints
type Controller struct {
	generator *service.LicenseGenerator
}

// NewController creates a new JetBrains controller
func NewController() *Controller {
	return &Controller{
		generator: service.NewLicenseGenerator(),
	}
}

// GenerateLicense handles license generation requests. ShouldBind dispatches
// on Content-Type: application/json hits the json tags, form-encoded bodies
// hit the form tags. The required+min=1 binding on LicenseeName covers the
// empty-name validation, so no manual check is needed.
func (c *Controller) GenerateLicense(ctx *gin.Context) {
	var req types.GenerateLicenseRequest
	if err := ctx.ShouldBind(&req); err != nil {
		httpx.HandleError(ctx, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	logger.Info("Generating JetBrains license for: " + req.LicenseeName)

	response, err := c.generator.GenerateLicense(req)
	if err != nil {
		logger.Error("Failed to generate license", err)
		httpx.HandleError(ctx, http.StatusInternalServerError, "Failed to generate license")
		return
	}

	httpx.HandleSuccess(ctx, response)
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

	httpx.HandleSuccess(ctx, powerConfig)
}

// FetchProductsLatest fetches the latest products
func (c *Controller) FetchProductsLatest(ctx *gin.Context) {
	go func() {
		if err := service.FetchLatestProducts(); err != nil {
			logger.Error("Failed to fetch latest products", err)
		}
	}()

	httpx.HandleSuccess(ctx, gin.H{
		"message": "Fetching latest products in background",
		"status":  "processing",
	})
}

// FetchPluginsLatest fetches the latest plugins
func (c *Controller) FetchPluginsLatest(ctx *gin.Context) {
	go func() {
		if err := service.FetchLatestPlugins(); err != nil {
			logger.Error("Failed to fetch latest plugins", err)
		}
	}()

	httpx.HandleSuccess(ctx, gin.H{
		"message": "Fetching latest plugins in background",
		"status":  "processing",
	})
}

// GetProducts returns all available products
func (c *Controller) GetProducts(ctx *gin.Context) {
	products, err := service.GetAllProducts()
	if err != nil {
		logger.Error("Failed to get products", err)
		httpx.HandleError(ctx, http.StatusInternalServerError, "Failed to get products")
		return
	}
	httpx.HandleSuccess(ctx, products)
}

// GetPlugins returns all available plugins
func (c *Controller) GetPlugins(ctx *gin.Context) {
	plugins, err := service.GetAllPlugins()
	if err != nil {
		logger.Error("Failed to get plugins", err)
		httpx.HandleError(ctx, http.StatusInternalServerError, "Failed to get plugins")
		return
	}
	httpx.HandleSuccess(ctx, plugins)
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
	generator *service.LicenseGenerator
}

// NewServerController creates a new server controller
func NewServerController() *ServerController {
	return &ServerController{
		generator: service.NewLicenseGenerator(),
	}
}

// LicenseServerRule returns the license server rules
func (sc *ServerController) LicenseServerRule(ctx *gin.Context) {
	powerConfig := sc.generator.GetPowerConfig()
	ctx.String(http.StatusOK, powerConfig.FullConfig)
}
