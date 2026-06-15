package api

import (
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	"license/internal/finalshell/service"
	"license/internal/httpx"

	"github.com/gin-gonic/gin"
)

// Controller defines the controller structure with performance monitoring
type Controller struct {
	requestCount int64
	totalTime    int64 // nanoseconds
	cacheHits    int64
	cacheMisses  int64
}

// NewController creates a new controller instance
func NewController() *Controller {
	return &Controller{}
}

// GenerateLicense generates a license handling function with optimizations
func (controller *Controller) GenerateLicense(c *gin.Context) {
	start := time.Now()

	machineCode := c.PostForm("machineCode")
	machineCode = strings.TrimSpace(machineCode)

	if machineCode == "" {
		httpx.HandleError(c, http.StatusBadRequest, "Machine code cannot be empty")
		return
	}

	// 输入验证 - 提前验证避免不必要的计算
	if len(machineCode) < 3 || len(machineCode) > 200 {
		httpx.HandleError(c, http.StatusBadRequest, "Invalid machine code format")
		return
	}

	// 检查是否来自缓存
	cacheSize := service.GetCacheStats()
	isCached := cacheSize > 0

	licenses := service.GenerateLicense(machineCode)

	// 添加缓存头
	c.Header("Cache-Control", "public, max-age=1800") // 30分钟
	if isCached {
		c.Header("X-Cache", "HIT")
		atomic.AddInt64(&controller.cacheHits, 1)
	} else {
		c.Header("X-Cache", "MISS")
		atomic.AddInt64(&controller.cacheMisses, 1)
	}

	// 性能监控
	duration := time.Since(start)
	atomic.AddInt64(&controller.totalTime, duration.Nanoseconds())
	atomic.AddInt64(&controller.requestCount, 1)

	// 在响应头中添加性能信息（开发环境）
	if gin.Mode() == gin.DebugMode {
		c.Header("X-Response-Time", duration.String())
		c.Header("X-Cache-Size", string(rune(cacheSize)))
	}

	httpx.HandleSuccess(c, licenses)
}

// GetStats returns performance statistics
func (controller *Controller) GetStats(c *gin.Context) {
	requestCount := atomic.LoadInt64(&controller.requestCount)
	totalTime := atomic.LoadInt64(&controller.totalTime)
	cacheHits := atomic.LoadInt64(&controller.cacheHits)
	cacheMisses := atomic.LoadInt64(&controller.cacheMisses)
	cacheSize := service.GetCacheStats()

	var avgTime int64
	if requestCount > 0 {
		avgTime = totalTime / requestCount
	}

	var hitRate float64
	if cacheHits+cacheMisses > 0 {
		hitRate = float64(cacheHits) / float64(cacheHits+cacheMisses) * 100
	}

	stats := gin.H{
		"total_requests":   requestCount,
		"average_time_ns":  avgTime,
		"average_time_ms":  float64(avgTime) / 1_000_000,
		"cache_size":       cacheSize,
		"cache_hits":       cacheHits,
		"cache_misses":     cacheMisses,
		"cache_hit_rate_%": hitRate,
	}
	httpx.HandleSuccess(c, stats)
}

// ClearCache clears the license cache
func (controller *Controller) ClearCache(c *gin.Context) {
	service.ClearCache()

	// Reset stats
	atomic.StoreInt64(&controller.requestCount, 0)
	atomic.StoreInt64(&controller.totalTime, 0)
	atomic.StoreInt64(&controller.cacheHits, 0)
	atomic.StoreInt64(&controller.cacheMisses, 0)

	httpx.HandleSuccess(c, gin.H{"message": "Cache cleared successfully"})
}
