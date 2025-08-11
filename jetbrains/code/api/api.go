package api

import (
	"license/jetbrains/code/service"
	"license/jetbrains/util"
	"license/logger"
	"strings"

	"github.com/gin-gonic/gin"
)

// Controller defines the controller structure
type Controller struct {
}

// NewController creates a new controller instance
func NewController() *Controller {
	return &Controller{}
}

// FetchProduceLatest fetches the latest activation codes
func (controller *Controller) FetchProduceLatest(c *gin.Context) {
	go func() {
		productService := service.NewProductService()
		err := productService.FetchLatest()
		if err != nil {
			logger.Error("Failed to fetch latest product:", err)
			return
		}
	}()

	c.JSON(200, gin.H{
		"message": "Fetching latest products in background",
	})
}

// FetchPluginLatest fetches the latest plugins
func (controller *Controller) FetchPluginLatest(c *gin.Context) {
	go func() {
		pluginService := service.NewPluginService()
		err := pluginService.FetchLatest()
		if err != nil {
			logger.Error("Failed to fetch latest plugin:", err)
			return
		}
	}()

	c.JSON(200, gin.H{
		"message": "Fetching latest plugins in background",
	})
}

// Generate generates activation codes
func (controller *Controller) Generate(c *gin.Context) {
	licenseeName := c.Query("licenseeName")
	effectiveDate := c.Query("effectiveDate")
	codes := c.Query("codes")

	// Split string into array
	var codesArray []string
	if len(codes) > 0 {
		codesArray = strings.Split(codes, ",")
	}

	// Generate license
	activationCode, err := service.GenerateLicense(licenseeName, effectiveDate, codesArray)
	if err != nil {
		logger.Error("Failed to generate license:", err)
		c.String(500, "Failed to generate license")
	}
	// Generate powerConf
	fakeCert := util.GetFake()
	powerConfRule := util.GeneratePowerResult(fakeCert.CodeCert, fakeCert.CodeRootCert)

	// Assemble data
	var result strings.Builder
	result.WriteString("================== power.conf ==================")
	result.WriteString("\n[Result]")
	result.WriteString("\n; Lemon active by code\n")
	result.WriteString(powerConfRule)
	result.WriteString("\n================== power.conf ==================")
	result.WriteString("\n================== activation code ==================\n")
	result.WriteString(activationCode)
	result.WriteString("\n================== activation code ==================\n")

	c.String(200, result.String())
}
