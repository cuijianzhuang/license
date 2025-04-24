package main

import (
	"embed"
	"flag"
	"fmt"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"license/config"
	"license/cron"
	"license/initialize"
	"license/logger"
	"license/router"
	"license/sys"
	"net/http"
	"strings"
)

//go:embed web/build
var EmbedFrontendFS embed.FS

func main() {
	// Define version flag
	versionFlag := flag.Bool("version", false, "Print version information and exit")
	flag.Parse()

	versionInfo := fmt.Sprintf("License v%s build(%s) %s", sys.Version, sys.Hash, sys.Arch)

	// If the version flag is specified, output version information from the constant and exit immediately
	if *versionFlag {
		// Print version information
		fmt.Println(versionInfo)
		return
	}

	// output version information to log
	logger.Sys(versionInfo)

	// Initialize global configuration
	config.InitConfig()
	// Initialize components
	initialize.ExecuteInitialize()
	// Initialize database
	config.SetupDatabase()
	// Initialize scheduled tasks
	cron.InitCron()

	// Set up GIN router
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// Set up routes under the /api prefix
	apiGroup := r.Group("/api")
	router.SetupRouter(apiGroup)

	// Create custom middleware to intercept all API requests
	r.Use(func(c *gin.Context) {
		path := c.Request.URL.Path

		// Ignore requests already starting with /api, these are handled by apiGroup
		if strings.HasPrefix(path, "/api/") {
			c.Next()
			return
		}

		// Check if it's an API request, if so redirect to the same name API handler
		if router.IsAPIPath(path) {
			// Copy the current request context, but pass the path to the API handler
			c.Request.URL.Path = path
			router.HandleAPIRequest(c)
			c.Abort() // Stop executing subsequent middleware
			return
		}

		// Not an API request, continue executing subsequent middleware
		c.Next()
	})

	// Serve frontend static files
	embedFS, err := static.EmbedFolder(EmbedFrontendFS, "web/build")
	if err != nil {
		logger.Error("Failed to load embedded frontend files", err)
		return
	}
	r.Use(static.Serve("/", embedFS))

	// Handle SPA routes
	r.NoRoute(func(c *gin.Context) {
		// If it's an API request, return 404 directly
		if len(c.Request.URL.Path) >= 4 && c.Request.URL.Path[:4] == "/api" {
			c.JSON(http.StatusNotFound, gin.H{"code": "PAGE_NOT_FOUND", "message": "API endpoint not found"})
			return
		}

		// For non-API requests, provide index.html to support SPA frontend routing
		c.Request.URL.Path = "/"
		r.HandleContext(c)
	})

	server := fmt.Sprintf("%s:%d", config.GetConfig().HttpHost, config.GetConfig().HttpPort)
	logger.Sys(fmt.Sprintf("Server starting, http://%s", server))
	// Start the server
	err = r.Run(server)
	if err != nil {
		logger.Error("Server failed to start", err)
		return
	}
}
