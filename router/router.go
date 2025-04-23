package router

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"license/config"
	finalshell "license/finalshell/api"
	gitlab "license/gitlab/api"
	jetbrainCode "license/jetbrains/code/api"
	jetbrainServer "license/jetbrains/server/api"
	jrebel "license/jrebel/api"
	"license/logger"
	mobaxterm "license/mobaxterm/api"
	rpc "license/rpc/controller"
	"license/utils/useragent"
	"net/http"
	"strings"
	"time"
)

// VersionResponse 定义版本响应结构
type VersionResponse struct {
	Version       string `json:"version"`
	NeedUpdate    bool   `json:"needUpdate"`
	LatestVersion string `json:"latestVersion,omitempty"`
}

// GitHubRelease GitHub API 版本发布响应结构
type GitHubRelease struct {
	TagName string `json:"tag_name"`
}

// 缓存GitHub最新版本的信息
var (
	cachedLatestVersion string
	lastFetchTime       time.Time
	cacheExpiration     = 30 * time.Minute
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

// 获取GitHub最新版本号（移除v前缀）
func getLatestVersionFromGitHub() string {
	// 如果缓存尚未过期，则使用缓存的值
	if !lastFetchTime.IsZero() && time.Since(lastFetchTime) < cacheExpiration && cachedLatestVersion != "" {
		return cachedLatestVersion
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", "https://api.github.com/repos/nannanStrawberry314/license/releases/latest", nil)
	if err != nil {
		logger.Error("创建GitHub API请求失败", err)
		return ""
	}

	// 使用随机UA
	req.Header.Set("User-Agent", useragent.GetRandom())
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil {
		logger.Error("请求GitHub API失败", err)
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Error("GitHub API返回非200状态码", nil)
		return ""
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("读取GitHub API响应失败", err)
		return ""
	}

	var release GitHubRelease
	if err := json.Unmarshal(body, &release); err != nil {
		logger.Error("解析GitHub API响应失败", err)
		return ""
	}

	// 去掉版本号前的"v"
	version := release.TagName
	if strings.HasPrefix(version, "v") {
		version = version[1:]
	}

	// 更新缓存
	cachedLatestVersion = version
	lastFetchTime = time.Now()

	return version
}

// 比较版本号大小
func compareVersions(current, latest string) bool {
	// 将两个版本按照点号分割
	currentParts := strings.Split(current, ".")
	latestParts := strings.Split(latest, ".")

	// 逐个比较版本号的各个部分
	for i := 0; i < len(currentParts) && i < len(latestParts); i++ {
		if currentParts[i] < latestParts[i] {
			return true // 需要更新
		} else if currentParts[i] > latestParts[i] {
			return false // 不需要更新
		}
	}

	// 如果前面的部分都相等，但是latest的部分更多，则需要更新
	return len(latestParts) > len(currentParts)
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
			currentVersion := config.Version
			latestVersion := getLatestVersionFromGitHub()
			needUpdate := false

			if latestVersion != "" {
				needUpdate = compareVersions(currentVersion, latestVersion)
			}

			c.JSON(200, VersionResponse{
				Version:       currentVersion,
				NeedUpdate:    needUpdate,
				LatestVersion: latestVersion,
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
