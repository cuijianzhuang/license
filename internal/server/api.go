package server

import (
	"license/internal/httpx"
	"license/internal/sys"

	"github.com/gin-gonic/gin"
)

// Controller Package server provides the server controller for handling API requests.
type Controller struct {
}

// NewServerController creates a new instance of the server controller.
func NewServerController() *Controller {
	return &Controller{}
}

// GetStatus get api status info
func (ctrl *Controller) GetStatus(c *gin.Context) {
	httpx.HandleSuccess(c, gin.H{
		"status": true,
	})
}

// GetVersion get api version info
func (ctrl *Controller) GetVersion(c *gin.Context) {
	currentVersion := sys.GetVersion()
	latestVersion := getLatestVersionFromGitHub()
	needUpdate := false

	if latestVersion != "" {
		needUpdate = compareVersions(currentVersion, latestVersion)
	}

	httpx.HandleSuccess(c, VersionResponse{
		Version:       currentVersion,
		Build:         sys.GetBuild(),
		OsArch:        sys.GetOsArch(),
		NeedUpdate:    needUpdate,
		LatestVersion: latestVersion,
	})
}
