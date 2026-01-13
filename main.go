package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"license/config"
	"license/cron"
	"license/initialize"
	"license/logger"
	"license/router"
	"license/sys"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

//go:embed web/build
var EmbedFrontendFS embed.FS

// 预加载的 index.html 内容
var indexHTMLContent []byte

// 静态资源扩展名集合（用于缓存判断）
var staticExtensions = map[string]bool{
	".js": true, ".css": true, ".png": true, ".jpg": true, ".jpeg": true,
	".gif": true, ".svg": true, ".ico": true, ".woff": true, ".woff2": true,
	".ttf": true, ".eot": true, ".map": true, ".webp": true,
}

// isStaticAsset 判断是否为静态资源
func isStaticAsset(path string) bool {
	if strings.Contains(path, "/static/") {
		return true
	}
	return staticExtensions[filepath.Ext(path)]
}

// setNoCacheHeaders 设置禁用缓存的响应头
func setNoCacheHeaders(c *gin.Context) {
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")
}

// serveIndexHTML 返回预加载的 index.html
func serveIndexHTML(c *gin.Context) {
	setNoCacheHeaders(c)
	c.Data(http.StatusOK, "text/html; charset=utf-8", indexHTMLContent)
}

// noRouteHandler handles 404s for SPA support
func noRouteHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		// API 路径返回 404 JSON
		if strings.HasPrefix(path, "/api") {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    "PAGE_NOT_FOUND",
				"message": "API endpoint not found",
			})
			return
		}

		// SPA 路由返回 index.html
		serveIndexHTML(c)
	}
}

// staticCacheMiddleware 静态资源缓存中间件
func staticCacheMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		if isStaticAsset(path) {
			// 静态资源缓存 1 年
			c.Header("Cache-Control", "public, max-age=31536000, immutable")
		} else if path == "/" || path == "/index.html" {
			setNoCacheHeaders(c)
		}

		c.Next()
	}
}

// setupStaticFileServer configures static file serving
func setupStaticFileServer(r *gin.Engine) error {
	// 获取 web/build 子文件系统
	webFS, err := fs.Sub(EmbedFrontendFS, "web/build")
	if err != nil {
		return fmt.Errorf("failed to get sub filesystem: %w", err)
	}

	// 添加缓存中间件
	r.Use(staticCacheMiddleware())

	// 预创建 FileServer
	fileServer := http.FileServer(http.FS(webFS))

	// 静态文件服务中间件
	r.Use(func(c *gin.Context) {
		path := c.Request.URL.Path

		// 跳过 API 请求
		if strings.HasPrefix(path, "/api") {
			c.Next()
			return
		}

		// 尝试打开文件（去掉开头的 /）
		filePath := strings.TrimPrefix(path, "/")
		if filePath == "" {
			filePath = "index.html"
		}

		file, err := webFS.Open(filePath)
		if err == nil {
			_ = file.Close()
			// 文件存在，使用 FileServer 处理
			fileServer.ServeHTTP(c.Writer, c.Request)
			c.Abort()
			return
		}

		// 文件不存在，继续到 NoRoute handler（返回 index.html 支持 SPA 路由）
		c.Next()
	})

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

	// 预加载 index.html
	var err error
	indexHTMLContent, err = EmbedFrontendFS.ReadFile("web/build/index.html")
	if err != nil {
		logger.Error("Failed to preload index.html: %v", err)
		return
	}

	// Set up GIN router
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	// Add essential middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Set up API routes under /api prefix
	apiGroup := r.Group("/api")
	router.SetupRouter(apiGroup)

	// Register external tool routes at root path (no conflict with frontend)
	// Includes: /rpc/*, /agent/*, /server/*
	rootGroup := r.Group("")
	router.SetupExternalRoutes(rootGroup)

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
