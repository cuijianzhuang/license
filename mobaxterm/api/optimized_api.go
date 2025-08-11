package api

import (
	"archive/zip"
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"license/logger"
	"license/mobaxterm/entity"
	"license/utils/useragent"
	v1 "license/v1"
	"net/http"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// 优化配置常量
const (
	licenseCacheMaxSize   = 5000            // 最大缓存许可证数量
	licenseCacheTTL       = 10 * time.Minute // 许可证缓存TTL
	versionCacheMaxAge    = 5 * time.Minute  // 版本缓存时间
	zipBufferPoolSize     = 100              // ZIP缓冲池大小
	cleanupInterval       = 2 * time.Minute  // 清理间隔
	maxConcurrency        = 150              // 最大并发数
)

// LicenseCacheEntry 许可证缓存条目
type LicenseCacheEntry struct {
	zipData   []byte
	timestamp time.Time
}

// LicenseCache 许可证缓存管理器
type LicenseCache struct {
	cache   sync.Map // 使用sync.Map保证并发安全
	entries int64    // 原子计数器
}

// VersionCache 版本缓存
type VersionCache struct {
	versions      []string
	lastFetchTime time.Time
	mutex         sync.RWMutex
}

// PerformanceStats 性能统计
type PerformanceStats struct {
	TotalRequests        int64     `json:"total_requests"`
	LicenseCacheHits     int64     `json:"license_cache_hits"`
	LicenseCacheMisses   int64     `json:"license_cache_misses"`
	VersionCacheHits     int64     `json:"version_cache_hits"`
	VersionCacheMisses   int64     `json:"version_cache_misses"`
	AverageGenTime       int64     `json:"average_gen_time_ns"`
	TotalGenTime         int64     `json:"total_gen_time_ns"`
	LicenseCacheHitRate  float64   `json:"license_cache_hit_rate"`
	VersionCacheHitRate  float64   `json:"version_cache_hit_rate"`
	LastCleanup          time.Time `json:"last_cleanup"`
	CurrentCacheSize     int64     `json:"current_cache_size"`
	PoolObjectsInUse     int64     `json:"pool_objects_in_use"`
	MemoryStats          MemoryStats `json:"memory_stats"`
	LoadLevel            string    `json:"load_level"`
}

// MemoryStats 内存使用统计
type MemoryStats struct {
	HeapAlloc     uint64    `json:"heap_alloc"`
	HeapSys       uint64    `json:"heap_sys"`
	HeapIdle      uint64    `json:"heap_idle"`
	HeapInuse     uint64    `json:"heap_inuse"`
	NumGoroutines int       `json:"num_goroutines"`
	NumGC         uint32    `json:"num_gc"`
	GCPauseTotal  uint64    `json:"gc_pause_total"`
	LastGCTime    time.Time `json:"last_gc_time"`
}

// OptimizedController 优化版MobaXterm控制器
type OptimizedController struct {
	licenseCache  *LicenseCache
	versionCache  *VersionCache
	stats         *PerformanceStats
	mu            sync.RWMutex

	// 对象池
	zipBufferPool     sync.Pool
	bytesBufferPool   sync.Pool
	licenseEntityPool sync.Pool
	stringBuilderPool sync.Pool

	// 清理任务控制
	cleanupTicker   *time.Ticker
	cleanupStopChan chan struct{}

	// 限流器
	requestLimiter chan struct{}
	maxConcurrency int

	// Base64编码表（预计算）
	variantBase64Map map[int]byte
}

// 全局实例
var (
	optimizedController     *OptimizedController
	optimizedControllerOnce sync.Once
)

// NewOptimizedController 创建优化版控制器实例
func NewOptimizedController() (*OptimizedController, error) {
	var initError error
	optimizedControllerOnce.Do(func() {
		controller := &OptimizedController{
			licenseCache:    &LicenseCache{},
			versionCache:    &VersionCache{},
			stats:           &PerformanceStats{LastCleanup: time.Now()},
			cleanupStopChan: make(chan struct{}),
			requestLimiter:  make(chan struct{}, maxConcurrency),
			maxConcurrency:  maxConcurrency,
		}

		// 初始化预计算的Base64映射表
		controller.initBase64Map()

		// 初始化对象池
		controller.initObjectPools()

		// 启动清理任务
		controller.startCleanupTask()

		optimizedController = controller
	})

	return optimizedController, initError
}

// initBase64Map 初始化Base64编码映射表
func (oc *OptimizedController) initBase64Map() {
	variantBase64Table := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/="
	oc.variantBase64Map = make(map[int]byte)
	for i, v := range variantBase64Table {
		oc.variantBase64Map[i] = byte(v)
	}
}

// initObjectPools 初始化对象池
func (oc *OptimizedController) initObjectPools() {
	// ZIP缓冲池
	oc.zipBufferPool = sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer(make([]byte, 0, 1024))
		},
	}

	// 字节缓冲池
	oc.bytesBufferPool = sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer(make([]byte, 0, 256))
		},
	}

	// 许可证实体对象池
	oc.licenseEntityPool = sync.Pool{
		New: func() interface{} {
			return &entity.License{}
		},
	}

	// 字符串构建器池
	oc.stringBuilderPool = sync.Pool{
		New: func() interface{} {
			builder := &strings.Builder{}
			builder.Grow(128)
			return builder
		},
	}
}

// startCleanupTask 启动定期清理任务
func (oc *OptimizedController) startCleanupTask() {
	oc.cleanupTicker = time.NewTicker(cleanupInterval)
	go func() {
		for {
			select {
			case <-oc.cleanupTicker.C:
				oc.cleanupExpiredCache()
			case <-oc.cleanupStopChan:
				return
			}
		}
	}()
}

// cleanupExpiredCache 清理过期缓存
func (oc *OptimizedController) cleanupExpiredCache() {
	now := time.Now()
	cleanedCount := 0

	oc.licenseCache.cache.Range(func(key, value interface{}) bool {
		entry := value.(*LicenseCacheEntry)
		if now.Sub(entry.timestamp) > licenseCacheTTL {
			oc.licenseCache.cache.Delete(key)
			atomic.AddInt64(&oc.licenseCache.entries, -1)
			cleanedCount++
		}
		return true
	})

	// 如果缓存过大，删除最旧的条目
	if atomic.LoadInt64(&oc.licenseCache.entries) > licenseCacheMaxSize {
		oc.evictOldestEntries(licenseCacheMaxSize / 4)
	}

	oc.mu.Lock()
	oc.stats.LastCleanup = now
	oc.stats.CurrentCacheSize = atomic.LoadInt64(&oc.licenseCache.entries)
	oc.mu.Unlock()
}

// evictOldestEntries 淘汰最旧的缓存条目
func (oc *OptimizedController) evictOldestEntries(count int64) {
	type keyTime struct {
		key  interface{}
		time time.Time
	}

	var entries []keyTime
	oc.licenseCache.cache.Range(func(key, value interface{}) bool {
		entry := value.(*LicenseCacheEntry)
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
		oc.licenseCache.cache.Delete(entries[i].key)
		atomic.AddInt64(&oc.licenseCache.entries, -1)
	}
}

// OptimizedFetchVersions 优化版获取版本号
func (oc *OptimizedController) OptimizedFetchVersions(c *gin.Context) {
	oc.versionCache.mutex.RLock()
	// 检查缓存是否有效
	if !oc.versionCache.lastFetchTime.IsZero() &&
		time.Since(oc.versionCache.lastFetchTime) < versionCacheMaxAge &&
		len(oc.versionCache.versions) > 0 {
		
		atomic.AddInt64(&oc.stats.VersionCacheHits, 1)
		versions := oc.versionCache.versions
		oc.versionCache.mutex.RUnlock()
		v1.HandleSuccess(c, versions)
		return
	}
	oc.versionCache.mutex.RUnlock()

	// 缓存未命中，获取新版本
	atomic.AddInt64(&oc.stats.VersionCacheMisses, 1)

	versions, err := oc.fetchVersionsFromWeb()
	if err != nil {
		logger.Error("", fmt.Errorf("failed to fetch versions: %v", err))
		// 返回默认版本列表
		defaultVersions := []string{"25.1", "25.0", "24.4", "24.3", "23.6", "23.5", "23.0"}
		v1.HandleSuccess(c, defaultVersions)
		return
	}

	// 更新缓存
	oc.versionCache.mutex.Lock()
	oc.versionCache.versions = versions
	oc.versionCache.lastFetchTime = time.Now()
	oc.versionCache.mutex.Unlock()

	v1.HandleSuccess(c, versions)
}

// fetchVersionsFromWeb 从网站获取版本信息
func (oc *OptimizedController) fetchVersionsFromWeb() ([]string, error) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest("GET", "https://mobaxterm.mobatek.net/download-home-edition.html", nil)
	if err != nil {
		return nil, err
	}

	headers := useragent.GetRandomWithAcceptHeaders()
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	versionsMap := make(map[string]bool)
	versionRegex := regexp.MustCompile(`Version (\d+\.\d+)`)

	doc.Find("p.version_titre").Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		matches := versionRegex.FindAllStringSubmatch(text, -1)
		for _, match := range matches {
			if len(match) >= 2 {
				versionsMap[match[1]] = true
			}
		}
	})

	versions := make([]string, 0, len(versionsMap))
	for version := range versionsMap {
		versions = append(versions, version)
	}

	if len(versions) == 0 {
		return []string{"25.1", "25.0", "24.4", "24.3", "23.6", "23.5", "23.0"}, nil
	}

	// 排序版本
	sort.Slice(versions, func(i, j int) bool {
		partsI := strings.Split(versions[i], ".")
		partsJ := strings.Split(versions[j], ".")

		majorI, _ := strconv.Atoi(partsI[0])
		majorJ, _ := strconv.Atoi(partsJ[0])

		if majorI != majorJ {
			return majorI > majorJ
		}

		minorI, _ := strconv.Atoi(partsI[1])
		minorJ, _ := strconv.Atoi(partsJ[1])

		return minorI > minorJ
	})

	return versions, nil
}

// OptimizedGenerateLicense 优化版生成许可证
func (oc *OptimizedController) OptimizedGenerateLicense(c *gin.Context) {
	start := time.Now()

	name := c.Query("name")
	if name == "" {
		name = c.PostForm("name")
	}
	version := c.Query("version")
	if version == "" {
		version = c.PostForm("version")
	}
	countStr := c.Query("count")
	if countStr == "" {
		countStr = c.PostForm("count")
	}

	count, err := strconv.Atoi(countStr)
	if err != nil {
		v1.HandleError(c, 400, "Invalid count value")
		return
	}

	if name == "" || version == "" || count <= 0 {
		v1.HandleError(c, 400, "Invalid parameters")
		return
	}

	// 生成缓存键
	cacheKey := fmt.Sprintf("%s:%s:%d", name, version, count)

	// 检查缓存
	if cachedEntry, found := oc.licenseCache.cache.Load(cacheKey); found {
		entry := cachedEntry.(*LicenseCacheEntry)
		if time.Since(entry.timestamp) < licenseCacheTTL {
			atomic.AddInt64(&oc.stats.LicenseCacheHits, 1)
			atomic.AddInt64(&oc.stats.TotalRequests, 1)
			
			c.Header("Content-Type", "application/zip")
			c.Header("Content-Disposition", "attachment; filename=Custom.mxtpro")
			c.Data(http.StatusOK, "application/zip", entry.zipData)
			return
		}
		// 过期条目，删除
		oc.licenseCache.cache.Delete(cacheKey)
		atomic.AddInt64(&oc.licenseCache.entries, -1)
	}

	// 缓存未命中，生成许可证
	atomic.AddInt64(&oc.stats.LicenseCacheMisses, 1)

	versionArr := strings.Split(version, ".")
	if len(versionArr) != 2 {
		v1.HandleError(c, 400, "Invalid version format")
		return
	}

	major, err := strconv.ParseInt(versionArr[0], 10, 64)
	if err != nil {
		v1.HandleError(c, 400, "Invalid major version")
		return
	}

	minor, err := strconv.ParseInt(versionArr[1], 10, 64)
	if err != nil {
		v1.HandleError(c, 400, "Invalid minor version")
		return
	}

	// 使用优化的许可证生成方法
	zipData, err := oc.generateOptimizedLicense(1, count, name, major, minor)
	if err != nil {
		logger.Error("", fmt.Errorf("failed to generate license: %v", err))
		v1.HandleError(c, 500, "Failed to generate license")
		return
	}

	// 添加到缓存
	if atomic.LoadInt64(&oc.licenseCache.entries) < licenseCacheMaxSize {
		oc.licenseCache.cache.Store(cacheKey, &LicenseCacheEntry{
			zipData:   zipData,
			timestamp: time.Now(),
		})
		atomic.AddInt64(&oc.licenseCache.entries, 1)
	}

	// 更新统计信息
	genTime := time.Since(start).Nanoseconds()
	atomic.AddInt64(&oc.stats.TotalGenTime, genTime)
	atomic.AddInt64(&oc.stats.TotalRequests, 1)

	c.Header("Content-Type", "application/zip")
	c.Header("Content-Disposition", "attachment; filename=Custom.mxtpro")
	c.Data(http.StatusOK, "application/zip", zipData)
}

// generateOptimizedLicense 优化版许可证生成
func (oc *OptimizedController) generateOptimizedLicense(userType, count int, username string, major, minor int64) ([]byte, error) {
	// 从对象池获取字符串构建器
	builder := oc.stringBuilderPool.Get().(*strings.Builder)
	defer func() {
		builder.Reset()
		oc.stringBuilderPool.Put(builder)
	}()

	// 高效构建许可证字符串
	builder.WriteString(strconv.Itoa(userType))
	builder.WriteByte('#')
	builder.WriteString(username)
	builder.WriteByte('|')
	builder.WriteString(strconv.FormatInt(major, 10))
	builder.WriteString(strconv.FormatInt(minor, 10))
	builder.WriteByte('#')
	builder.WriteString(strconv.Itoa(count))
	builder.WriteByte('#')
	builder.WriteString(strconv.FormatInt(major, 10))
	builder.WriteByte('3')
	builder.WriteString(strconv.FormatInt(minor, 10))
	builder.WriteByte('6')
	builder.WriteString(strconv.FormatInt(minor, 10))
	builder.WriteString("#0#0#0#")

	licenseString := builder.String()

	// 使用优化的加密和编码
	encryptedBytes := oc.optimizedEncryptBytes(0x787, []byte(licenseString))
	encodedBytes := oc.optimizedVariantBase64Encode(encryptedBytes)

	// 使用对象池获取ZIP缓冲区
	zipBuffer := oc.zipBufferPool.Get().(*bytes.Buffer)
	defer func() {
		zipBuffer.Reset()
		oc.zipBufferPool.Put(zipBuffer)
	}()

	// 创建ZIP文件
	zipWriter := zip.NewWriter(zipBuffer)
	
	header := &zip.FileHeader{
		Name:               "Pro.key",
		Method:             zip.Store,
		CompressedSize64:   uint64(len(encodedBytes)),
		UncompressedSize64: uint64(len(encodedBytes)),
	}
	header.SetModTime(time.Now())

	writer, err := zipWriter.CreateRaw(header)
	if err != nil {
		return nil, err
	}

	_, err = writer.Write(encodedBytes)
	if err != nil {
		return nil, err
	}

	err = zipWriter.Close()
	if err != nil {
		return nil, err
	}

	// 复制数据以避免池对象被回收时数据丢失
	result := make([]byte, zipBuffer.Len())
	copy(result, zipBuffer.Bytes())

	return result, nil
}

// optimizedEncryptBytes 优化版字节加密
func (oc *OptimizedController) optimizedEncryptBytes(key int, bs []byte) []byte {
	result := make([]byte, 0, len(bs))
	for _, b := range bs {
		encryptedByte := b ^ byte((key>>8)&0xff)
		result = append(result, encryptedByte)
		key = (int(result[len(result)-1]) & key) | 0x482D
	}
	return result
}

// optimizedVariantBase64Encode 优化版Base64编码
func (oc *OptimizedController) optimizedVariantBase64Encode(bs []byte) []byte {
	result := make([]byte, 0, (len(bs)*4+2)/3)
	blocksCount := len(bs) / 3
	leftBytes := len(bs) % 3

	for i := 0; i < blocksCount; i++ {
		codingInt := oc.littleEndianBytes(bs[3*i : 3*i+3])
		result = append(result,
			oc.variantBase64Map[codingInt&0x3f],
			oc.variantBase64Map[(codingInt>>6)&0x3f],
			oc.variantBase64Map[(codingInt>>12)&0x3f],
			oc.variantBase64Map[(codingInt>>18)&0x3f],
		)
	}

	if leftBytes == 1 {
		codingInt := oc.littleEndianBytes(bs[3*blocksCount:])
		result = append(result,
			oc.variantBase64Map[codingInt&0x3f],
			oc.variantBase64Map[(codingInt>>6)&0x3f],
		)
	} else if leftBytes == 2 {
		codingInt := oc.littleEndianBytes(bs[3*blocksCount:])
		result = append(result,
			oc.variantBase64Map[codingInt&0x3f],
			oc.variantBase64Map[(codingInt>>6)&0x3f],
			oc.variantBase64Map[(codingInt>>12)&0x3f],
		)
	}

	return result
}

// littleEndianBytes 小端字节转换
func (oc *OptimizedController) littleEndianBytes(bs []byte) int {
	result := int(bs[0])
	for i := 1; i < len(bs); i++ {
		result = result | int(bs[i])<<(8*i)
	}
	return result
}

// GetPerformanceStats 获取性能统计信息
func (oc *OptimizedController) GetPerformanceStats(c *gin.Context) {
	oc.mu.RLock()
	defer oc.mu.RUnlock()

	stats := *oc.stats
	stats.CurrentCacheSize = atomic.LoadInt64(&oc.licenseCache.entries)

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
	stats.LoadLevel = oc.calculateLoadLevel(stats)

	totalRequests := atomic.LoadInt64(&oc.stats.TotalRequests)
	if totalRequests > 0 {
		stats.AverageGenTime = atomic.LoadInt64(&oc.stats.TotalGenTime) / totalRequests

		licenseHits := atomic.LoadInt64(&oc.stats.LicenseCacheHits)
		licenseMisses := atomic.LoadInt64(&oc.stats.LicenseCacheMisses)
		if licenseHits+licenseMisses > 0 {
			stats.LicenseCacheHitRate = float64(licenseHits) / float64(licenseHits+licenseMisses) * 100
		}

		versionHits := atomic.LoadInt64(&oc.stats.VersionCacheHits)
		versionMisses := atomic.LoadInt64(&oc.stats.VersionCacheMisses)
		if versionHits+versionMisses > 0 {
			stats.VersionCacheHitRate = float64(versionHits) / float64(versionHits+versionMisses) * 100
		}
	}

	c.JSON(http.StatusOK, stats)
}

// calculateLoadLevel 计算当前负载级别
func (oc *OptimizedController) calculateLoadLevel(stats PerformanceStats) string {
	memoryUsageRate := float64(stats.MemoryStats.HeapInuse) / float64(stats.MemoryStats.HeapSys)
	goroutineCount := stats.MemoryStats.NumGoroutines
	cacheUsageRate := float64(stats.CurrentCacheSize) / float64(licenseCacheMaxSize)

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
func (oc *OptimizedController) HealthCheck(c *gin.Context) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	loadLevel := oc.calculateLoadLevel(PerformanceStats{
		MemoryStats: MemoryStats{
			HeapInuse:     memStats.HeapInuse,
			HeapSys:       memStats.HeapSys,
			NumGoroutines: runtime.NumGoroutine(),
		},
		CurrentCacheSize: atomic.LoadInt64(&oc.licenseCache.entries),
	})

	health := gin.H{
		"status":          "ok",
		"load_level":      loadLevel,
		"memory_mb":       memStats.HeapAlloc / 1024 / 1024,
		"goroutines":      runtime.NumGoroutine(),
		"license_cache":   atomic.LoadInt64(&oc.licenseCache.entries),
		"uptime":          time.Since(oc.stats.LastCleanup),
	}

	switch loadLevel {
	case "critical":
		c.JSON(http.StatusServiceUnavailable, health)
	case "high":
		c.JSON(http.StatusTooManyRequests, health)
	default:
		c.JSON(http.StatusOK, health)
	}
}

// ClearCache 清空缓存
func (oc *OptimizedController) ClearCache(c *gin.Context) {
	cleared := int64(0)
	oc.licenseCache.cache.Range(func(key, value interface{}) bool {
		oc.licenseCache.cache.Delete(key)
		cleared++
		return true
	})

	atomic.StoreInt64(&oc.licenseCache.entries, 0)

	// 重置统计信息
	atomic.StoreInt64(&oc.stats.LicenseCacheHits, 0)
	atomic.StoreInt64(&oc.stats.LicenseCacheMisses, 0)
	atomic.StoreInt64(&oc.stats.VersionCacheHits, 0)
	atomic.StoreInt64(&oc.stats.VersionCacheMisses, 0)
	atomic.StoreInt64(&oc.stats.TotalRequests, 0)
	atomic.StoreInt64(&oc.stats.TotalGenTime, 0)

	// 清空版本缓存
	oc.versionCache.mutex.Lock()
	oc.versionCache.versions = nil
	oc.versionCache.lastFetchTime = time.Time{}
	oc.versionCache.mutex.Unlock()

	c.JSON(http.StatusOK, gin.H{
		"message":         "Cache cleared successfully",
		"cleared_entries": cleared,
	})
}

// RateLimitMiddleware 限流中间件
func (oc *OptimizedController) RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		select {
		case oc.requestLimiter <- struct{}{}:
			defer func() { <-oc.requestLimiter }()
			c.Next()
		default:
			var memStats runtime.MemStats
			runtime.ReadMemStats(&memStats)

			loadLevel := "critical"
			if runtime.NumGoroutine() < 500 {
				loadLevel = "high"
			}

			c.Header("Retry-After", "5")
			c.Header("X-Rate-Limit-Limit", fmt.Sprintf("%d", oc.maxConcurrency))
			c.Header("X-Rate-Limit-Remaining", "0")
			c.Header("X-Load-Level", loadLevel)

			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":           "Rate limit exceeded",
				"message":         "Server is currently under high load, please retry after 5 seconds",
				"load_level":      loadLevel,
				"max_concurrency": oc.maxConcurrency,
			})
			c.Abort()
		}
	}
}

// ForceGC 强制垃圾回收
func (oc *OptimizedController) ForceGC(c *gin.Context) {
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

// Shutdown 优雅关闭
func (oc *OptimizedController) Shutdown() {
	if oc.cleanupTicker != nil {
		oc.cleanupTicker.Stop()
	}
	if oc.cleanupStopChan != nil {
		close(oc.cleanupStopChan)
	}
}