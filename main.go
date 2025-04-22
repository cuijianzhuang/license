package main

import (
	"embed"
	"fmt"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"license/config"
	"license/cron"
	"license/initialize"
	"license/logger"
	"license/router"
	"net/http"
)

//go:embed web/build
var EmbedFrontendFS embed.FS

func main() {
	// 初始化全局配置
	config.InitConfig()
	// 初始化
	initialize.ExecuteInitialize()
	// 初始化数据库
	config.SetupDatabase()
	// 初始化定时任务
	cron.InitCron()

	// 设置 GIN 路由
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// 为API路由添加前缀
	apiGroup := r.Group("/api")
	router.SetupRouter(apiGroup)

	// 提供前端静态文件
	embedFS, err := static.EmbedFolder(EmbedFrontendFS, "web/build")
	if err != nil {
		logger.Error("加载嵌入式前端文件失败", err)
		return
	}
	r.Use(static.Serve("/", embedFS))

	// 处理SPA路由
	r.NoRoute(func(c *gin.Context) {
		// 如果是API请求，直接返回404
		if len(c.Request.URL.Path) >= 4 && c.Request.URL.Path[:4] == "/api" {
			c.JSON(http.StatusNotFound, gin.H{"code": "PAGE_NOT_FOUND", "message": "API endpoint not found"})
			return
		}

		// 对于非API请求，提供index.html以支持SPA前端路由
		c.Request.URL.Path = "/"
		r.HandleContext(c)
	})

	server := fmt.Sprintf("%s:%d", config.GetConfig().HttpHost, config.GetConfig().HttpPort)
	logger.Sys(fmt.Sprintf("服务启动中, http://%s", server))
	// 启动服务器
	err = r.Run(server)
	if err != nil {
		logger.Error("服务器启动失败", err)
		return
	}
}
