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

	// jrebel - 使用优化版本
	jrebelLeasesApi := jrebel.NewLeasesController()  // 原始版本备用
	jrebelIndexApi := jrebel.NewIndexController()
	
	// 创建优化版控制器
	optimizedJrebelApi, err := jrebel.NewOptimizedLeasesController()
	var useOptimized bool
	if err != nil {
		// 如果优化版本初始化失败，回退到原始版本
		useOptimized = false
	} else {
		useOptimized = true
	}
	
	jrebelGroup := r.Group("/jrebel")
	{
		jrebelGroup.GET("/", jrebelIndexApi.IndexHandler)
		
		if useOptimized {
			// 使用优化版本
			jrebelGroup.DELETE("/leases/1", optimizedJrebelApi.OptimizedLeases1Handler)
			jrebelGroup.POST("/leases", optimizedJrebelApi.OptimizedLeasesHandler)
			jrebelGroup.POST("/validate-connection", optimizedJrebelApi.OptimizedValidateHandler)
			jrebelGroup.POST("/features", optimizedJrebelApi.OptimizedValidateHandler)
			jrebelGroup.GET("/features", optimizedJrebelApi.OptimizedValidateHandler)
			// 优化版本专有端点
			jrebelGroup.GET("/performance-stats", optimizedJrebelApi.GetPerformanceStats)
			jrebelGroup.POST("/clear-cache", optimizedJrebelApi.ClearCache)
		} else {
			// 回退到原始版本
			jrebelGroup.DELETE("/leases/1", jrebelLeasesApi.Leases1Handler)
			jrebelGroup.POST("/leases", jrebelLeasesApi.LeasesHandler)
			jrebelGroup.POST("/validate-connection", jrebelLeasesApi.ValidateHandler)
			jrebelGroup.POST("/features", jrebelLeasesApi.ValidateHandler)
			jrebelGroup.GET("/features", jrebelLeasesApi.ValidateHandler)
		}
		
		jrebelGroup.POST("/leases/1", func(c *gin.Context) {
			c.Status(405)
		})
	}
	jrebelAgentGroup := r.Group("/agent")
	{
		if useOptimized {
			jrebelAgentGroup.DELETE("/leases/1", optimizedJrebelApi.OptimizedLeases1Handler)
			jrebelAgentGroup.POST("/leases", optimizedJrebelApi.OptimizedLeasesHandler)
			jrebelAgentGroup.POST("/validate-connection", optimizedJrebelApi.OptimizedValidateHandler)
			jrebelAgentGroup.POST("/features", optimizedJrebelApi.OptimizedValidateHandler)
			jrebelAgentGroup.GET("/features", optimizedJrebelApi.OptimizedValidateHandler)
		} else {
			jrebelAgentGroup.DELETE("/leases/1", jrebelLeasesApi.Leases1Handler)
			jrebelAgentGroup.POST("/leases", jrebelLeasesApi.LeasesHandler)
			jrebelAgentGroup.POST("/validate-connection", jrebelLeasesApi.ValidateHandler)
			jrebelAgentGroup.POST("/features", jrebelLeasesApi.ValidateHandler)
			jrebelAgentGroup.GET("/features", jrebelLeasesApi.ValidateHandler)
		}
		
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
