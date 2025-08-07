package router

import (
	"github.com/gin-gonic/gin"
	finalshell "license/finalshell/api"
	gitlab "license/gitlab/api"
	jetbrainCode "license/jetbrains/api"
	jetbrainServer "license/jetbrains/api"
	jrebel "license/jrebel/api"
	mobaxterm "license/mobaxterm/api"
	rpc "license/rpc/controller"
	"license/server"
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
	serverApi := server.NewServerController()
	serverGroup := r.Group("/server")
	{
		serverGroup.GET("/status", serverApi.GetStatus)
		serverGroup.GET("/version", serverApi.GetVersion)
	}

	// final-shell
	finalShellApi := finalshell.NewController()
	finalShellGroup := r.Group("/final-shell")
	{
		finalShellGroup.POST("/generateLicense", finalShellApi.GenerateLicense)
		finalShellGroup.GET("/stats", finalShellApi.GetStats)
		finalShellGroup.POST("/clear-cache", finalShellApi.ClearCache)
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
	jetbrainsServerApi := jetbrainServer.NewServerController()
	jetbrainsCodeApi := jetbrainCode.NewController()

	jetbrainsGroup := r.Group("/jetbrains")
	{
		// License generation
		jetbrainsGroup.GET("/generate", jetbrainsCodeApi.GenerateLicense)
		jetbrainsGroup.POST("/generate", jetbrainsCodeApi.GenerateLicense)
		
		// Power config
		jetbrainsGroup.GET("/licenseServerRule", jetbrainsServerApi.LicenseServerRule)
		jetbrainsGroup.GET("/powerConfig", jetbrainsCodeApi.GetPowerConfig)
		
		// Product and plugin management
		jetbrainsGroup.GET("/products", jetbrainsCodeApi.GetProducts)
		jetbrainsGroup.GET("/products/fetchLatest", jetbrainsCodeApi.FetchProductsLatest)
		jetbrainsGroup.GET("/plugins", jetbrainsCodeApi.GetPlugins)
		jetbrainsGroup.GET("/plugins/fetchLatest", jetbrainsCodeApi.FetchPluginsLatest)
		
		// Health check
		jetbrainsGroup.GET("/health", jetbrainsCodeApi.HealthCheck)
		
		// Backward compatibility
		jetbrainsGroup.GET("/product/fetchLatest", jetbrainsCodeApi.FetchProductsLatest)
		jetbrainsGroup.GET("/plugin/fetchLatest", jetbrainsCodeApi.FetchPluginsLatest)
	}
}
