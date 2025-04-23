package router

import (
	"github.com/gin-gonic/gin"
	"license/config"
	finalshell "license/finalshell/api"
	gitlab "license/gitlab/api"
	jetbrainCode "license/jetbrains/code/api"
	jetbrainServer "license/jetbrains/server/api"
	jrebel "license/jrebel/api"
	mobaxterm "license/mobaxterm/api"
	rpc "license/rpc/controller"
	"strings"
)

// List of API path prefixes
var apiPrefixes = []string{
	"/server/",
	"/final-shell/",
	"/gitlab/",
	"/rpc/",
	"/jrebel/",
	"/agent/",
	"/mobaxterm/",
	"/jetbrains/",
}

// IsAPIPath determines if the given path is an API path
func IsAPIPath(path string) bool {
	for _, prefix := range apiPrefixes {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}
	return false
}

// HandleAPIRequest handles API requests
func HandleAPIRequest(c *gin.Context) {
	// Create a temporary routing engine to handle the request
	tmpEngine := gin.New()
	tmpGroup := tmpEngine.Group("/")

	// Set up the router
	SetupRouter(tmpGroup)

	// Handle the request
	tmpEngine.HandleContext(c)
}

func SetupRouter(r *gin.RouterGroup) {
	serverGroup := r.Group("/server")
	{
		serverGroup.GET("/status", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status": true,
			})
		})
		
		serverGroup.GET("/version", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"version": config.GetVersion(),
			})
		})
	}

	// final-shell
	finalShellApi := finalshell.NewController()
	finalShellGroup := r.Group("/final-shell")
	{
		finalShellGroup.POST("/generateLicense", finalShellApi.GenerateLicense)
	}

	// gitlab
	gitlabApi := gitlab.NewController()
	gitlabGroup := r.Group("/gitlab")
	{
		gitlabGroup.POST("/generate", gitlabApi.Generate)
	}

	// rpc
	rpcApi := rpc.NewRpcController()
	rpcGroup := r.Group("/rpc")
	{
		rpcGroup.GET("/ping.action", rpcApi.Ping)
		rpcGroup.GET("/obtainTicket.action", rpcApi.ObtainTicket)
		rpcGroup.GET("/releaseTicket.action", rpcApi.ReleaseTicket)
	}

	// jrebel
	jrebelLeasesApi := jrebel.NewLeasesController()
	jrebelIndexApi := jrebel.NewIndexController()
	jrebelGroup := r.Group("/jrebel")
	{
		jrebelGroup.GET("/", jrebelIndexApi.IndexHandler)
		jrebelGroup.DELETE("/leases/1", jrebelLeasesApi.Leases1Handler)
		jrebelGroup.POST("/leases/1", func(c *gin.Context) {
			c.Status(405)
		})
		jrebelGroup.POST("/leases", jrebelLeasesApi.LeasesHandler)
		jrebelGroup.POST("/validate-connection", jrebelLeasesApi.ValidateHandler)
		jrebelGroup.POST("/features", jrebelLeasesApi.ValidateHandler)
		jrebelGroup.GET("/features", jrebelLeasesApi.ValidateHandler)
	}
	jrebelAgentGroup := r.Group("/agent")
	{
		jrebelAgentGroup.DELETE("/leases/1", jrebelLeasesApi.Leases1Handler)
		jrebelAgentGroup.POST("/leases/1", func(c *gin.Context) {
			c.Status(405)
		})
		jrebelAgentGroup.POST("/leases", jrebelLeasesApi.LeasesHandler)
		jrebelAgentGroup.POST("/validate-connection", jrebelLeasesApi.ValidateHandler)
		jrebelAgentGroup.POST("/features", jrebelLeasesApi.ValidateHandler)
		jrebelAgentGroup.GET("/features", jrebelLeasesApi.ValidateHandler)
	}

	// mobaxterm
	mobaxtermApi := mobaxterm.NewMobaXtermController()
	mobaxtermGroup := r.Group("/mobaxterm")
	{
		mobaxtermGroup.POST("/generate", mobaxtermApi.GenerateLicense)
		mobaxtermGroup.GET("/versions", func(c *gin.Context) {
			// Add no-cache headers
			c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
			c.Header("Pragma", "no-cache")
			c.Header("Expires", "0")
			
			mobaxtermApi.FetchVersions(c)
		})
	}

	// jetbrains
	jetbrainsServerApi := jetbrainServer.NewLicenseServerController()
	jetbrainsCodeApi := jetbrainCode.NewController()

	jetbrainsGroup := r.Group("/jetbrains")
	{
		jetbrainsGroup.GET("/licenseServerRule", jetbrainsServerApi.LicenseServerRule)
		jetbrainsGroup.GET("/product/fetchLatest", jetbrainsCodeApi.FetchProduceLatest)
		jetbrainsGroup.GET("/plugin/fetchLatest", jetbrainsCodeApi.FetchPluginLatest)
		jetbrainsGroup.GET("/generate", jetbrainsCodeApi.Generate)
	}
}
