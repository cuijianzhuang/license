package api

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"net/http"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// 缓存配置
const (
	signatureCacheMaxSize = 2000            // 最大缓存签名数量
	signatureCacheTTL     = 5 * time.Minute // 缓存TTL
	objectPoolMaxSize     = 500             // 对象池最大大小
	cleanupInterval       = 1 * time.Minute // 缓存清理间隔
)

// SignatureCacheEntry 签名缓存条目
type SignatureCacheEntry struct {
	signature string
	timestamp time.Time
}

// SignatureCache 签名缓存管理器
type SignatureCache struct {
	cache   sync.Map // 使用sync.Map保证并发安全
	entries int64    // 原子计数器
}

// PerformanceStats 性能统计
type PerformanceStats struct {
	TotalRequests    int64       `json:"total_requests"`
	CacheHits        int64       `json:"cache_hits"`
	CacheMisses      int64       `json:"cache_misses"`
	AverageSignTime  int64       `json:"average_sign_time_ns"`
	TotalSignTime    int64       `json:"total_sign_time_ns"`
	CacheHitRate     float64     `json:"cache_hit_rate"`
	LastCleanup      time.Time   `json:"last_cleanup"`
	CurrentCacheSize int64       `json:"current_cache_size"`
	PoolObjectsInUse int64       `json:"pool_objects_in_use"`
	MemoryStats      MemoryStats `json:"memory_stats"`
	LoadLevel        string      `json:"load_level"` // "low", "medium", "high", "critical"
}

// MemoryStats 内存使用统计
type MemoryStats struct {
	HeapAlloc     uint64    `json:"heap_alloc"`     // 当前堆内存分配
	HeapSys       uint64    `json:"heap_sys"`       // 系统堆内存
	HeapIdle      uint64    `json:"heap_idle"`      // 空闲堆内存
	HeapInuse     uint64    `json:"heap_inuse"`     // 使用中堆内存
	NumGoroutines int       `json:"num_goroutines"` // Goroutine数量
	NumGC         uint32    `json:"num_gc"`         // GC次数
	GCPauseTotal  uint64    `json:"gc_pause_total"` // GC总暂停时间
	LastGCTime    time.Time `json:"last_gc_time"`   // 上次GC时间
}

// LeasesController 优化版JRebel控制器 (合并后的版本)
type LeasesController struct {
	privateKey     *rsa.PrivateKey   // 预解析的私钥
	signatureCache *SignatureCache   // 签名缓存
	stats          *PerformanceStats // 性能统计
	mu             sync.RWMutex      // 读写锁
}

// 全局实例
var (
	controllerInstance     *LeasesController
	controllerInstanceOnce sync.Once
)

// NewLeasesController 创建控制器实例 (使用优化版本)
func NewLeasesController() (*LeasesController, error) {
	var initError error
	controllerInstanceOnce.Do(func() {
		controller := &LeasesController{
			signatureCache: &SignatureCache{},
			stats:          &PerformanceStats{LastCleanup: time.Now()},
		}

		// 预解析RSA私钥
		if err := controller.initPrivateKey(); err != nil {
			initError = err
			return
		}

		// 启动清理任务
		controller.startCleanupTask()

		controllerInstance = controller
	})

	if initError != nil {
		return nil, initError
	}

	return controllerInstance, nil
}

// initPrivateKey 初始化RSA私钥（一次性解析）
func (lc *LeasesController) initPrivateKey() error {
	block, _ := pem.Decode([]byte(LeasesPrivateKey))
	if block == nil {
		return fmt.Errorf("failed to decode PEM block containing the private key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse RSA private key: %v", err)
	}

	lc.privateKey = privateKey
	return nil
}

// startCleanupTask 启动定期清理任务
func (lc *LeasesController) startCleanupTask() {
	ticker := time.NewTicker(cleanupInterval)
	go func() {
		for range ticker.C {
			lc.cleanupExpiredCache()
		}
	}()
}

// cleanupExpiredCache 清理过期缓存
func (lc *LeasesController) cleanupExpiredCache() {
	now := time.Now()
	cleanedCount := 0

	lc.signatureCache.cache.Range(func(key, value interface{}) bool {
		entry := value.(*SignatureCacheEntry)
		if now.Sub(entry.timestamp) > signatureCacheTTL {
			lc.signatureCache.cache.Delete(key)
			atomic.AddInt64(&lc.signatureCache.entries, -1)
			cleanedCount++
		}
		return true
	})

	// 如果缓存过大，删除最旧的条目
	if atomic.LoadInt64(&lc.signatureCache.entries) > signatureCacheMaxSize {
		lc.evictOldestEntries(signatureCacheMaxSize / 4) // 删除1/4的条目
	}

	lc.mu.Lock()
	lc.stats.LastCleanup = now
	lc.stats.CurrentCacheSize = atomic.LoadInt64(&lc.signatureCache.entries)
	lc.mu.Unlock()
}

// evictOldestEntries 淘汰最旧的缓存条目 (优化版排序算法)
func (lc *LeasesController) evictOldestEntries(count int64) {
	type keyTime struct {
		key  interface{}
		time time.Time
	}

	var entries []keyTime
	lc.signatureCache.cache.Range(func(key, value interface{}) bool {
		entry := value.(*SignatureCacheEntry)
		entries = append(entries, keyTime{key: key, time: entry.timestamp})
		return true
	})

	// 使用Go标准库的sort.Slice进行高效排序
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].time.Before(entries[j].time)
	})

	// 删除最旧的count个条目
	deleteCount := count
	if int64(len(entries)) < deleteCount {
		deleteCount = int64(len(entries))
	}

	for i := int64(0); i < deleteCount; i++ {
		lc.signatureCache.cache.Delete(entries[i].key)
		atomic.AddInt64(&lc.signatureCache.entries, -1)
	}
}

// sign 优化版签名函数
func (lc *LeasesController) sign(clientRandomness, guid string, offline bool, validFrom, validUntil int64) string {
	start := time.Now()

	// 生成缓存键
	cacheKey := fmt.Sprintf("%s:%s:%t:%d:%d", clientRandomness, guid, offline, validFrom, validUntil)

	// 检查缓存
	if cachedEntry, found := lc.signatureCache.cache.Load(cacheKey); found {
		entry := cachedEntry.(*SignatureCacheEntry)
		if time.Since(entry.timestamp) < signatureCacheTTL {
			atomic.AddInt64(&lc.stats.CacheHits, 1)
			return entry.signature
		}
		// 过期条目，删除
		lc.signatureCache.cache.Delete(cacheKey)
		atomic.AddInt64(&lc.signatureCache.entries, -1)
	}

	// 缓存未命中，计算签名
	atomic.AddInt64(&lc.stats.CacheMisses, 1)

	var builder strings.Builder
	builder.Grow(128)
	builder.WriteString(clientRandomness)
	builder.WriteByte(';')
	builder.WriteString(ServerRandomness)
	builder.WriteByte(';')
	builder.WriteString(guid)
	builder.WriteByte(';')
	builder.WriteString(strconv.FormatBool(offline))

	if offline {
		builder.WriteByte(';')
		builder.WriteString(strconv.FormatInt(validFrom, 10))
		builder.WriteByte(';')
		builder.WriteString(strconv.FormatInt(validUntil, 10))
	}

	hasher := sha1.New()
	hasher.Write([]byte(builder.String()))
	hashed := hasher.Sum(nil)

	// RSA签名（使用预解析的私钥）
	signature, err := rsa.SignPKCS1v15(rand.Reader, lc.privateKey, crypto.SHA1, hashed)
	if err != nil {
		return ""
	}

	signatureStr := base64.StdEncoding.EncodeToString(signature)

	// 添加到缓存
	if atomic.LoadInt64(&lc.signatureCache.entries) < signatureCacheMaxSize {
		lc.signatureCache.cache.Store(cacheKey, &SignatureCacheEntry{
			signature: signatureStr,
			timestamp: time.Now(),
		})
		atomic.AddInt64(&lc.signatureCache.entries, 1)
	}

	// 更新统计信息
	signTime := time.Since(start).Nanoseconds()
	atomic.AddInt64(&lc.stats.TotalSignTime, signTime)
	atomic.AddInt64(&lc.stats.TotalRequests, 1)

	return signatureStr
}

// LeasesHandler handles the "/leases" endpoint (优化版).
func (lc *LeasesController) LeasesHandler(c *gin.Context) {
	clientRandomness := c.PostForm("randomness")
	username := c.PostForm("username")
	guid := c.PostForm("guid")
	offline, _ := strconv.ParseBool(c.PostForm("offline"))
	clientTime, _ := strconv.ParseInt(c.PostForm("clientTime"), 10, 64)

	var validFrom, validUntil int64
	if offline {
		expiration := clientTime + 180*24*60*60*1000
		validFrom = clientTime
		validUntil = expiration
	}

	signature := lc.sign(clientRandomness, guid, offline, validFrom, validUntil)

	resp := LeasesResponse{
		ServerVersion:         ServerVersion,
		ServerProtocolVersion: ServerProtocolVersion,
		ServerGuid:            ServerGuid,
		GroupType:             GroupType,
		ID:                    1,
		LicenseType:           1,
		EvaluationLicense:     false,
		Signature:             signature,
		ServerRandomness:      ServerRandomness,
		SeatPoolType:          SeatPoolType,
		StatusCode:            StatusCode,
		Offline:               offline,
		ValidFrom:             validFrom,
		ValidUntil:            validUntil,
		Company:               username,
		OrderId:               uuid.NewString(),
		ZeroIds:               make([]string, 0),
		LicenseValidFrom:      1490544001000,
		LicenseValidUntil:     4102415999000,
	}

	c.Header("X-Cache", fmt.Sprintf("hits:%d,misses:%d",
		atomic.LoadInt64(&lc.stats.CacheHits),
		atomic.LoadInt64(&lc.stats.CacheMisses)))

	c.JSON(http.StatusOK, resp)
}

// Leases1Handler handles the "/leases/1" endpoint (优化版).
func (lc *LeasesController) Leases1Handler(c *gin.Context) {
	username := c.DefaultQuery("username", "")
	clientRandomness := c.PostForm("randomness")
	guid := c.PostForm("guid")
	offline, _ := strconv.ParseBool(c.PostForm("offline"))
	clientTime, _ := strconv.ParseInt(c.PostForm("clientTime"), 10, 64)

	var validFrom, validUntil int64
	if offline {
		expiration := clientTime + 180*24*60*60*1000
		validFrom = clientTime
		validUntil = expiration
	}

	signature := lc.sign(clientRandomness, guid, offline, validFrom, validUntil)

	resp := LeasesOneResponse{
		ServerVersion:         ServerVersion,
		ServerProtocolVersion: ServerProtocolVersion,
		ServerGuid:            ServerGuid,
		Signature:             signature,
		ServerRandomness:      ServerRandomness,
		Features:              "{}",
		GroupType:             GroupType,
		StatusCode:            StatusCode,
		Company:               username,
		Msg:                   "",
		StatusMessage:         "",
	}

	c.JSON(http.StatusOK, resp)
}

// ValidateHandler handles the "/validate-connection" endpoint (优化版).
func (lc *LeasesController) ValidateHandler(c *gin.Context) {
	clientRandomness := c.PostForm("randomness")
	guid := c.PostForm("guid")
	offline, _ := strconv.ParseBool(c.PostForm("offline"))
	clientTime, _ := strconv.ParseInt(c.PostForm("clientTime"), 10, 64)

	var validFrom, validUntil int64
	if offline {
		expiration := clientTime + 180*24*time.Hour.Milliseconds()
		validFrom = clientTime
		validUntil = expiration
	}

	signature := lc.sign(clientRandomness, guid, offline, validFrom, validUntil)

	resp := ValidateResponse{
		ServerVersion:         ServerVersion,
		ServerProtocolVersion: ServerProtocolVersion,
		ServerGuid:            ServerGuid,
		Signature:             signature,
		ServerRandomness:      ServerRandomness,
		Features:              "{}",
		GroupType:             GroupType,
		StatusCode:            StatusCode,
		Company:               Company,
		CanGetLease:           true,
		LicenseType:           "1",
		EvaluationLicense:     false,
		SeatPoolType:          SeatPoolType,
	}

	c.JSON(http.StatusOK, resp)
}

// GetPerformanceStats 获取性能统计信息
func (lc *LeasesController) GetPerformanceStats(c *gin.Context) {
	lc.mu.RLock()
	defer lc.mu.RUnlock()

	stats := *lc.stats
	stats.CurrentCacheSize = atomic.LoadInt64(&lc.signatureCache.entries)

	// 获取内存统计信息
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	stats.MemoryStats = MemoryStats{
		HeapAlloc:     memStats.HeapAlloc,
		HeapSys:       memStats.HeapSys,
		HeapIdle:      memStats.HeapIdle,
		HeapInuse:     memStats.HeapInuse,
		NumGoroutines: runtime.NumGoroutine(),
		NumGC:         memStats.NumGC,
		GCPauseTotal:  memStats.PauseTotalNs,
		LastGCTime:    time.Unix(0, int64(memStats.LastGC)),
	}

	// 计算负载级别
	stats.LoadLevel = lc.calculateLoadLevel(stats)

	totalRequests := atomic.LoadInt64(&lc.stats.TotalRequests)
	if totalRequests > 0 {
		stats.AverageSignTime = atomic.LoadInt64(&lc.stats.TotalSignTime) / totalRequests

		cacheHits := atomic.LoadInt64(&lc.stats.CacheHits)
		cacheMisses := atomic.LoadInt64(&lc.stats.CacheMisses)
		if cacheHits+cacheMisses > 0 {
			stats.CacheHitRate = float64(cacheHits) / float64(cacheHits+cacheMisses) * 100
		}
	}

	c.JSON(http.StatusOK, stats)
}

// calculateLoadLevel 计算当前负载级别
func (lc *LeasesController) calculateLoadLevel(stats PerformanceStats) string {
	// 基于内存使用率计算负载级别
	memoryUsageRate := float64(stats.MemoryStats.HeapInuse) / float64(stats.MemoryStats.HeapSys)

	// 基于Goroutine数量计算负载
	goroutineCount := stats.MemoryStats.NumGoroutines

	// 基于缓存大小计算负载
	cacheUsageRate := float64(stats.CurrentCacheSize) / float64(signatureCacheMaxSize)

	// 综合评估负载级别
	if memoryUsageRate > 0.9 || goroutineCount > 1000 || cacheUsageRate > 0.95 {
		return "critical"
	} else if memoryUsageRate > 0.7 || goroutineCount > 500 || cacheUsageRate > 0.8 {
		return "high"
	} else if memoryUsageRate > 0.5 || goroutineCount > 200 || cacheUsageRate > 0.6 {
		return "medium"
	}
	return "low"
}

// HealthCheck 健康检查端点
func (lc *LeasesController) HealthCheck(c *gin.Context) {
	stats := PerformanceStats{}
	stats.LoadLevel = lc.calculateLoadLevel(stats)

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	health := gin.H{
		"status":     "ok",
		"load_level": stats.LoadLevel,
		"memory_mb":  memStats.HeapAlloc / 1024 / 1024,
		"goroutines": runtime.NumGoroutine(),
		"cache_size": atomic.LoadInt64(&lc.signatureCache.entries),
		"uptime":     time.Since(lc.stats.LastCleanup),
	}

	// 根据负载级别设置HTTP状态码
	switch stats.LoadLevel {
	case "critical":
		c.JSON(http.StatusServiceUnavailable, health)
	case "high":
		c.JSON(http.StatusTooManyRequests, health)
	default:
		c.JSON(http.StatusOK, health)
	}
}

// ClearCache 清空缓存
func (lc *LeasesController) ClearCache(c *gin.Context) {
	cleared := int64(0)
	lc.signatureCache.cache.Range(func(key, value interface{}) bool {
		lc.signatureCache.cache.Delete(key)
		cleared++
		return true
	})

	atomic.StoreInt64(&lc.signatureCache.entries, 0)

	// 重置统计信息
	atomic.StoreInt64(&lc.stats.CacheHits, 0)
	atomic.StoreInt64(&lc.stats.CacheMisses, 0)
	atomic.StoreInt64(&lc.stats.TotalRequests, 0)
	atomic.StoreInt64(&lc.stats.TotalSignTime, 0)

	c.JSON(http.StatusOK, gin.H{
		"message":         "Cache cleared successfully",
		"cleared_entries": cleared,
	})
}

// ForceGC 强制垃圾回收（调试用）
func (lc *LeasesController) ForceGC(c *gin.Context) {
	var beforeStats, afterStats runtime.MemStats
	runtime.ReadMemStats(&beforeStats)

	runtime.GC()

	runtime.ReadMemStats(&afterStats)

	c.JSON(http.StatusOK, gin.H{
		"message":              "Garbage collection completed",
		"before_heap_alloc_mb": beforeStats.HeapAlloc / 1024 / 1024,
		"after_heap_alloc_mb":  afterStats.HeapAlloc / 1024 / 1024,
		"freed_mb":             (beforeStats.HeapAlloc - afterStats.HeapAlloc) / 1024 / 1024,
		"gc_count":             afterStats.NumGC,
	})
}
