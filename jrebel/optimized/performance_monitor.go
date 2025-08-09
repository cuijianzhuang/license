package optimized

import (
	"encoding/json"
	"fmt"
	"license/jrebel/constant"
	"sync"
	"sync/atomic"
	"time"
)

// PerformanceMonitor tracks performance metrics for the optimized JRebel service
type PerformanceMonitor struct {
	// Request counters
	totalRequests       int64
	cacheHits          int64
	cacheMisses        int64
	
	// Timing metrics
	totalResponseTime  int64 // in nanoseconds
	maxResponseTime    int64
	minResponseTime    int64
	
	// Error counters
	signatureErrors    int64
	parseErrors        int64
	
	// Memory pool statistics
	poolGets           int64
	poolPuts           int64
	poolMisses         int64
	
	// Concurrent access protection
	mutex              sync.RWMutex
	
	// Start time for uptime calculation
	startTime          time.Time
}

// PerformanceStats represents exportable performance statistics
type PerformanceStats struct {
	// Request statistics
	TotalRequests     int64   `json:"total_requests"`
	CacheHitRate      float64 `json:"cache_hit_rate_percent"`
	CacheHits         int64   `json:"cache_hits"`
	CacheMisses       int64   `json:"cache_misses"`
	
	// Timing statistics
	AverageResponseTime string `json:"average_response_time"`
	MaxResponseTime     string `json:"max_response_time"`
	MinResponseTime     string `json:"min_response_time"`
	
	// Error statistics
	SignatureErrors   int64   `json:"signature_errors"`
	ParseErrors       int64   `json:"parse_errors"`
	ErrorRate         float64 `json:"error_rate_percent"`
	
	// Memory pool statistics
	PoolEfficiency    float64 `json:"pool_efficiency_percent"`
	PoolGets          int64   `json:"pool_gets"`
	PoolPuts          int64   `json:"pool_puts"`
	PoolMisses        int64   `json:"pool_misses"`
	
	// System statistics
	Uptime            string `json:"uptime"`
	RequestsPerSecond float64 `json:"requests_per_second"`
}

// NewPerformanceMonitor creates a new performance monitor instance
func NewPerformanceMonitor() *PerformanceMonitor {
	return &PerformanceMonitor{
		startTime:       time.Now(),
		minResponseTime: ^int64(0), // Initialize to max value
	}
}

// RecordRequest records a successful request with its response time
func (pm *PerformanceMonitor) RecordRequest(responseTime time.Duration, cacheHit bool) {
	atomic.AddInt64(&pm.totalRequests, 1)
	responseTimeNanos := responseTime.Nanoseconds()
	atomic.AddInt64(&pm.totalResponseTime, responseTimeNanos)
	
	if cacheHit {
		atomic.AddInt64(&pm.cacheHits, 1)
	} else {
		atomic.AddInt64(&pm.cacheMisses, 1)
	}
	
	// Update min/max response times
	pm.mutex.Lock()
	if responseTimeNanos > pm.maxResponseTime {
		pm.maxResponseTime = responseTimeNanos
	}
	if responseTimeNanos < pm.minResponseTime {
		pm.minResponseTime = responseTimeNanos
	}
	pm.mutex.Unlock()
}

// RecordSignatureError records a signature generation error
func (pm *PerformanceMonitor) RecordSignatureError() {
	atomic.AddInt64(&pm.signatureErrors, 1)
}

// RecordParseError records a parsing error
func (pm *PerformanceMonitor) RecordParseError() {
	atomic.AddInt64(&pm.parseErrors, 1)
}

// RecordPoolOperation records object pool get/put operations
func (pm *PerformanceMonitor) RecordPoolGet() {
	atomic.AddInt64(&pm.poolGets, 1)
}

func (pm *PerformanceMonitor) RecordPoolPut() {
	atomic.AddInt64(&pm.poolPuts, 1)
}

func (pm *PerformanceMonitor) RecordPoolMiss() {
	atomic.AddInt64(&pm.poolMisses, 1)
}

// GetStats returns current performance statistics
func (pm *PerformanceMonitor) GetStats() PerformanceStats {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()
	
	totalRequests := atomic.LoadInt64(&pm.totalRequests)
	cacheHits := atomic.LoadInt64(&pm.cacheHits)
	cacheMisses := atomic.LoadInt64(&pm.cacheMisses)
	totalResponseTime := atomic.LoadInt64(&pm.totalResponseTime)
	signatureErrors := atomic.LoadInt64(&pm.signatureErrors)
	parseErrors := atomic.LoadInt64(&pm.parseErrors)
	poolGets := atomic.LoadInt64(&pm.poolGets)
	poolPuts := atomic.LoadInt64(&pm.poolPuts)
	poolMisses := atomic.LoadInt64(&pm.poolMisses)
	
	uptime := time.Since(pm.startTime)
	
	stats := PerformanceStats{
		TotalRequests:   totalRequests,
		CacheHits:       cacheHits,
		CacheMisses:     cacheMisses,
		SignatureErrors: signatureErrors,
		ParseErrors:     parseErrors,
		PoolGets:        poolGets,
		PoolPuts:        poolPuts,
		PoolMisses:      poolMisses,
		Uptime:         uptime.String(),
	}
	
	// Calculate cache hit rate
	totalCacheOps := cacheHits + cacheMisses
	if totalCacheOps > 0 {
		stats.CacheHitRate = float64(cacheHits) / float64(totalCacheOps) * 100
	}
	
	// Calculate average response time
	if totalRequests > 0 {
		avgNanos := totalResponseTime / totalRequests
		stats.AverageResponseTime = time.Duration(avgNanos).String()
		stats.RequestsPerSecond = float64(totalRequests) / uptime.Seconds()
	}
	
	// Format min/max response times
	if pm.maxResponseTime > 0 {
		stats.MaxResponseTime = time.Duration(pm.maxResponseTime).String()
	}
	if pm.minResponseTime < ^int64(0) {
		stats.MinResponseTime = time.Duration(pm.minResponseTime).String()
	}
	
	// Calculate error rate
	totalErrors := signatureErrors + parseErrors
	if totalRequests > 0 {
		stats.ErrorRate = float64(totalErrors) / float64(totalRequests) * 100
	}
	
	// Calculate pool efficiency
	if poolGets > 0 {
		stats.PoolEfficiency = float64(poolPuts) / float64(poolGets) * 100
	}
	
	return stats
}

// GetStatsJSON returns performance statistics as JSON
func (pm *PerformanceMonitor) GetStatsJSON() ([]byte, error) {
	stats := pm.GetStats()
	return json.MarshalIndent(stats, "", "  ")
}

// Reset resets all performance counters
func (pm *PerformanceMonitor) Reset() {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	
	atomic.StoreInt64(&pm.totalRequests, 0)
	atomic.StoreInt64(&pm.cacheHits, 0)
	atomic.StoreInt64(&pm.cacheMisses, 0)
	atomic.StoreInt64(&pm.totalResponseTime, 0)
	atomic.StoreInt64(&pm.signatureErrors, 0)
	atomic.StoreInt64(&pm.parseErrors, 0)
	atomic.StoreInt64(&pm.poolGets, 0)
	atomic.StoreInt64(&pm.poolPuts, 0)
	atomic.StoreInt64(&pm.poolMisses, 0)
	
	pm.maxResponseTime = 0
	pm.minResponseTime = ^int64(0)
	pm.startTime = time.Now()
}

// Enhanced OptimizedLeasesController with performance monitoring
type MonitoredOptimizedLeasesController struct {
	*OptimizedLeasesController
	monitor *PerformanceMonitor
}

// NewMonitoredOptimizedLeasesController creates a new monitored controller
func NewMonitoredOptimizedLeasesController() (*MonitoredOptimizedLeasesController, error) {
	controller, err := NewOptimizedLeasesController()
	if err != nil {
		return nil, err
	}
	
	return &MonitoredOptimizedLeasesController{
		OptimizedLeasesController: controller,
		monitor: NewPerformanceMonitor(),
	}, nil
}

// monitoredOptimizedSign wraps the optimized sign function with performance monitoring
func (controller *MonitoredOptimizedLeasesController) monitoredOptimizedSign(clientRandomness, guid string, offline bool, validFrom, validUntil int64) string {
	startTime := time.Now()
	
	// Check if this will be a cache hit
	cacheKey := fmt.Sprintf("%s:%s:%s:%t:%d:%d", 
		clientRandomness, constant.ServerRandomness, guid, offline, validFrom, validUntil)
	isHit := controller.signatureCache.get(cacheKey) != ""
	
	// Call the actual sign function
	signature := controller.optimizedSign(clientRandomness, guid, offline, validFrom, validUntil)
	
	// Record performance metrics
	responseTime := time.Since(startTime)
	controller.monitor.RecordRequest(responseTime, isHit)
	
	if signature == "" {
		controller.monitor.RecordSignatureError()
	}
	
	return signature
}

// GetPerformanceStats returns current performance statistics
func (controller *MonitoredOptimizedLeasesController) GetPerformanceStats() PerformanceStats {
	return controller.monitor.GetStats()
}

// GetPerformanceStatsJSON returns performance statistics as JSON
func (controller *MonitoredOptimizedLeasesController) GetPerformanceStatsJSON() ([]byte, error) {
	return controller.monitor.GetStatsJSON()
}

// ResetPerformanceStats resets all performance counters
func (controller *MonitoredOptimizedLeasesController) ResetPerformanceStats() {
	controller.monitor.Reset()
}
