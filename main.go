package main

import (
	"embed"
	"flag"
	"fmt"
	"license/config"
	"license/cron"
	"license/initialize"
	"license/logger"
	"license/router"
	"license/sys"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

//go:embed web/build
var EmbedFrontendFS embed.FS

// Cache for path checks to improve performance
var (
	pathCache = make(map[string]bool)
	cacheMu   sync.RWMutex
	cacheSize = 1000 // Maximum cache entries
)

// isAPIPathCached checks if a path is an API path with caching
func isAPIPathCached(path string) bool {
	cacheMu.RLock()
	if cached, exists := pathCache[path]; exists {
		cacheMu.RUnlock()
		return cached
	}
	cacheMu.RUnlock()

	// Check if it's an API path
	isAPI := router.IsAPIPath(path)

	// Cache the result if cache is not full
	cacheMu.Lock()
	if len(pathCache) < cacheSize {
		pathCache[path] = isAPI
	}
	cacheMu.Unlock()

	return isAPI
}

// apiMiddleware is an optimized middleware for API request handling
func apiMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		// Skip if already has /api/ prefix - these are handled by apiGroup
		if strings.HasPrefix(path, "/api/") {
			c.Next()
			return
		}

		// Use cached API path check
		if isAPIPathCached(path) {
			c.Request.URL.Path = path
			router.HandleAPIRequest(c)
			c.Abort()
			return
		}

		c.Next()
	}
}

// noRouteHandler handles 404s with better performance
func noRouteHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		// Fast check for API paths using string prefix
		if len(path) >= 4 && path[:4] == "/api" {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    "PAGE_NOT_FOUND",
				"message": "API endpoint not found",
			})
			return
		}

		// For non-API requests, serve SPA index.html
		// This ensures that all frontend routes (like /mobaxterm) serve the main app
		indexHTML, err := EmbedFrontendFS.ReadFile("web/build/index.html")
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to load index.html")
			return
		}
		
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")
		c.Data(http.StatusOK, "text/html; charset=utf-8", indexHTML)
	}
}

// setupStaticFileServer configures static file serving with caching
func setupStaticFileServer(r *gin.Engine) error {
	embedFS, err := static.EmbedFolder(EmbedFrontendFS, "web/build")
	if err != nil {
		return fmt.Errorf("failed to load embedded frontend files: %w", err)
	}

	// Add caching middleware for static files
	r.Use(func(c *gin.Context) {
		path := c.Request.URL.Path

		// Add cache headers for static assets
		if strings.Contains(path, "/static/") ||
			strings.HasSuffix(path, ".js") ||
			strings.HasSuffix(path, ".css") ||
			strings.HasSuffix(path, ".png") ||
			strings.HasSuffix(path, ".jpg") ||
			strings.HasSuffix(path, ".jpeg") ||
			strings.HasSuffix(path, ".gif") ||
			strings.HasSuffix(path, ".svg") ||
			strings.HasSuffix(path, ".ico") ||
			strings.HasSuffix(path, ".woff") ||
			strings.HasSuffix(path, ".woff2") ||
			strings.HasSuffix(path, ".ttf") ||
			strings.HasSuffix(path, ".eot") {

			// Cache static assets for 1 year
			c.Header("Cache-Control", "public, max-age=31536000, immutable")
			c.Header("Expires", time.Now().Add(365*24*time.Hour).Format(http.TimeFormat))
		} else if path == "/" || path == "/index.html" {
			// Don't cache the main HTML file to ensure updates are loaded
			c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
			c.Header("Pragma", "no-cache")
			c.Header("Expires", "0")
		}

		c.Next()
	})

	r.Use(static.Serve("/", embedFS))
	return nil
}

func main() {
	// Define version flag
	versionFlag := flag.Bool("version", false, "Print version information and exit")
	flag.Parse()

	versionInfo := fmt.Sprintf("License v%s build(%s) %s", sys.GetVersion(), sys.GetBuild(), sys.GetOsArch())

	// If the version flag is specified, output version information and exit
	if *versionFlag {
		fmt.Println(versionInfo)
		return
	}

	// Output version information to log
	logger.Sys(versionInfo)

	// Initialize global configuration
	config.InitConfig()

	// Initialize database
	config.SetupDatabase()

	// Initialize components
	if err := initialize.ExecuteInitialize(); err != nil {
		logger.Error("Failed to initialize components: %v", err)
		return
	}

	// Initialize scheduled tasks
	cron.InitCron()

	// Set up GIN router with optimized configuration
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	// Add essential middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Add compression middleware for better performance
	r.Use(func(c *gin.Context) {
		// Enable gzip compression for text responses
		if !strings.Contains(c.GetHeader("Accept-Encoding"), "gzip") {
			c.Next()
			return
		}

		// Set compression header for text-based content
		contentType := c.GetHeader("Content-Type")
		if strings.Contains(contentType, "application/json") ||
			strings.Contains(contentType, "text/") ||
			strings.Contains(contentType, "application/javascript") ||
			strings.Contains(contentType, "application/xml") {
			c.Header("Content-Encoding", "gzip")
		}

		c.Next()
	})

	// Set up API routes under the /api prefix
	apiGroup := r.Group("/api")
	router.SetupRouter(apiGroup)

	// Use middleware for API request interception
	r.Use(apiMiddleware())

	// Configure static file server with caching
	if err := setupStaticFileServer(r); err != nil {
		logger.Error("Failed to setup static file server", err)
		return
	}

	// Set NoRoute handler for SPA support
	r.NoRoute(noRouteHandler())

	// Start the server
	serverAddr := fmt.Sprintf("%s:%d", config.GetConfig().HttpHost, config.GetConfig().HttpPort)
	logger.Sys(fmt.Sprintf("Server starting, http://%s", serverAddr))

	if err := r.Run(serverAddr); err != nil {
		logger.Error("Server failed to start", err)
		return
	}
}
