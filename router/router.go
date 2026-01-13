package router

import (
	finalshell "license/finalshell/api"
	gitlab "license/gitlab/api"
	jetbrainServer "license/jetbrains/code/api/v2"
	jrebel "license/jrebel/api"
	mobaxterm "license/mobaxterm/api"
	rpc "license/rpc/controller"
	"license/server"

	"github.com/gin-gonic/gin"
)

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
	jrebelLeasesApi, _ := jrebel.NewLeasesController()
	jrebelIndexApi := jrebel.NewIndexController()

	jrebelGroup := r.Group("/jrebel")
	{
		jrebelGroup.GET("/", jrebelIndexApi.IndexHandler)

		// 使用优化版本（现在是默认版本）
		jrebelGroup.DELETE("/leases/1", jrebelLeasesApi.Leases1Handler)
		jrebelGroup.POST("/leases", jrebelLeasesApi.LeasesHandler)
		jrebelGroup.POST("/validate-connection", jrebelLeasesApi.ValidateHandler)
		jrebelGroup.POST("/features", jrebelLeasesApi.ValidateHandler)
		jrebelGroup.GET("/features", jrebelLeasesApi.ValidateHandler)
		// 优化版本专有端点
		jrebelGroup.GET("/performance-stats", jrebelLeasesApi.GetPerformanceStats)
		jrebelGroup.POST("/clear-cache", jrebelLeasesApi.ClearCache)
		jrebelGroup.GET("/health", jrebelLeasesApi.HealthCheck)
		jrebelGroup.POST("/force-gc", jrebelLeasesApi.ForceGC)

		jrebelGroup.POST("/leases/1", func(c *gin.Context) {
			c.Status(405)
		})
	}
	jrebelAgentGroup := r.Group("/agent")
	{
		jrebelAgentGroup.DELETE("/leases/1", jrebelLeasesApi.Leases1Handler)
		jrebelAgentGroup.POST("/leases", jrebelLeasesApi.LeasesHandler)
		jrebelAgentGroup.POST("/validate-connection", jrebelLeasesApi.ValidateHandler)
		jrebelAgentGroup.POST("/features", jrebelLeasesApi.ValidateHandler)
		jrebelAgentGroup.GET("/features", jrebelLeasesApi.ValidateHandler)

		jrebelAgentGroup.POST("/leases/1", func(c *gin.Context) {
			c.Status(405)
		})
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
		// 优化版本专有端点
		mobaxtermGroup.GET("/performance-stats", mobaxtermApi.GetPerformanceStats)
		mobaxtermGroup.POST("/clear-cache", mobaxtermApi.ClearCache)
		mobaxtermGroup.GET("/health", mobaxtermApi.HealthCheck)
		mobaxtermGroup.POST("/force-gc", mobaxtermApi.ForceGC)
	}

	// jetbrains
	jetbrainsServerApi := jetbrainServer.NewServerController()
	jetbrainsCodeApi := jetbrainServer.NewController()

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
