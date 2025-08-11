package api

import (
	"license/jetbrains/util"
	"strings"

	"github.com/gin-gonic/gin"
)

// LicenseServerController defines the controller structure
type LicenseServerController struct {
}

// NewLicenseServerController creates a new controller instance
func NewLicenseServerController() *LicenseServerController {
	return &LicenseServerController{}
}

// LicenseServerRule generates a license handling function
func (controller *LicenseServerController) LicenseServerRule(c *gin.Context) {
	fakeCert := util.GetFake()
	codePower := util.GeneratePowerResult(fakeCert.CodeCert, fakeCert.CodeRootCert)
	serverPower := util.GeneratePowerResult(fakeCert.ServerCert, fakeCert.ServerRootCert)

	// Construct the result
	var result strings.Builder
	result.WriteString("[Result]\n; Lemon active by code\n")
	result.WriteString(codePower)
	result.WriteString("\n[Result]\n; Lemon active by server\n")
	result.WriteString(serverPower)
	result.WriteString("\n")
	c.String(200, result.String())

}
