package server

import (
	"github.com/gin-gonic/gin"
	"license/sys"
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
	c.JSON(200, gin.H{
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

	c.JSON(200, VersionResponse{
		Version:       currentVersion,
		Hash:          sys.GetHash(),
		Arch:          sys.GetArch(),
		NeedUpdate:    needUpdate,
		LatestVersion: latestVersion,
	})
}
