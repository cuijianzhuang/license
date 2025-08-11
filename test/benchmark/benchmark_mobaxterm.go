//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"time"

	"license/mobaxterm/api"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("MobaXterm 许可证服务性能基准测试")
	fmt.Println("===============================")

	gin.SetMode(gin.ReleaseMode)

	// 创建控制器（现在已合并优化功能）
	controller := api.NewMobaXtermController()

	// 创建路由
	router := setupRouter(controller)

	// 运行基准测试
	runBenchmarkTests(router)
}

func setupRouter(controller *api.Controller) *gin.Engine {
	r := gin.New()
	// 应用限流中间件
	r.Use(controller.RateLimitMiddleware())
	
	r.GET("/mobaxterm/versions", controller.FetchVersions)
	r.POST("/mobaxterm/generate", controller.GenerateLicense)
	r.GET("/mobaxterm/generate", controller.GenerateLicense)
	r.GET("/mobaxterm/stats", controller.GetPerformanceStats)
	r.POST("/mobaxterm/clear-cache", controller.ClearCache)
	r.GET("/mobaxterm/health", controller.HealthCheck)
	r.POST("/mobaxterm/force-gc", controller.ForceGC)
	return r
}

func runBenchmarkTests(router *gin.Engine) {
	testCases := []struct {
		name     string
		endpoint string
		method   string
		testFunc func(*gin.Engine, string, string, int) time.Duration
	}{
		{"License Generation", "/mobaxterm/generate", "POST", benchmarkLicenseGeneration},
		{"Version Fetch", "/mobaxterm/versions", "GET", benchmarkVersionFetch},
	}

	for _, tc := range testCases {
		fmt.Printf("\n=== %s 性能测试 ===\n", tc.name)
		
		// 清理缓存
		if tc.name == "License Generation" {
			clearCache(router)
		}
		fmt.Println("无缓存测试:")
		noCacheTime := tc.testFunc(router, tc.endpoint, tc.method, 100)
		
		// 预热缓存
		if tc.name == "License Generation" {
			preloadLicenseCache(router, 10)
		}
		fmt.Println("缓存测试:")
		cachedTime := tc.testFunc(router, tc.endpoint, tc.method, 100)
		
		// 计算提升
		improvement := float64(noCacheTime) / float64(cachedTime)
		fmt.Printf("缓存优化倍数: %.2fx\n", improvement)
		
		// 显示缓存统计
		showCacheStats(router)
	}

	// 并发测试
	fmt.Println("\n=== 并发性能测试 ===")
	runConcurrencyTest(router)

	// 压力测试
	fmt.Println("\n=== 压力测试 ===")
	runStressTest(router)
}

func benchmarkLicenseGeneration(router *gin.Engine, endpoint, method string, iterations int) time.Duration {
	start := time.Now()
	
	for i := 0; i < iterations; i++ {
		req := createLicenseRequest(endpoint, method, i)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		if w.Code != http.StatusOK && w.Code != http.StatusTooManyRequests {
			fmt.Printf("Request failed with status: %d\n", w.Code)
		}
	}
	
	elapsed := time.Since(start)
	fmt.Printf("  %d 个请求用时: %v (平均: %v)\n", 
		iterations, elapsed, elapsed/time.Duration(iterations))
	
	return elapsed
}

func benchmarkVersionFetch(router *gin.Engine, endpoint, method string, iterations int) time.Duration {
	start := time.Now()
	
	for i := 0; i < iterations; i++ {
		req, _ := http.NewRequest(method, endpoint, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		if w.Code != http.StatusOK {
			fmt.Printf("Request failed with status: %d\n", w.Code)
		}
	}
	
	elapsed := time.Since(start)
	fmt.Printf("  %d 个请求用时: %v (平均: %v)\n", 
		iterations, elapsed, elapsed/time.Duration(iterations))
	
	return elapsed
}

func createLicenseRequest(endpoint, method string, seed int) *http.Request {
	form := url.Values{}
	form.Add("name", fmt.Sprintf("user-%d", seed))
	form.Add("version", "25.1")
	form.Add("count", "1")

	body := strings.NewReader(form.Encode())
	req, _ := http.NewRequest(method, endpoint, body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	
	return req
}

func preloadLicenseCache(router *gin.Engine, count int) {
	for i := 0; i < count; i++ {
		req := createLicenseRequest("/mobaxterm/generate", "POST", i)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}

func clearCache(router *gin.Engine) {
	req, _ := http.NewRequest("POST", "/mobaxterm/clear-cache", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
}

func showCacheStats(router *gin.Engine) {
	req, _ := http.NewRequest("GET", "/mobaxterm/stats", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	fmt.Printf("缓存统计: %s\n", w.Body.String())
}

func runConcurrencyTest(router *gin.Engine) {
	concurrencyLevels := []int{10, 50, 100}
	requestsPerLevel := 50 // 减少请求数以避免生成太多许可证文件
	
	for _, concurrency := range concurrencyLevels {
		fmt.Printf("\n并发数: %d\n", concurrency)
		
		clearCache(router)
		testTime := runConcurrentRequests(router, "/mobaxterm/generate", "POST", 
			concurrency, requestsPerLevel)
		
		fmt.Printf("用时: %v (QPS: %.2f)\n", testTime, 
			float64(concurrency*requestsPerLevel)/testTime.Seconds())
	}
}

func runConcurrentRequests(router *gin.Engine, endpoint, method string, 
	concurrency, requestsPerWorker int) time.Duration {
	
	var wg sync.WaitGroup
	start := time.Now()
	
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			
			for j := 0; j < requestsPerWorker; j++ {
				req := createLicenseRequest(endpoint, method, workerID*requestsPerWorker+j)
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)
			}
		}(i)
	}
	
	wg.Wait()
	return time.Since(start)
}

func runStressTest(router *gin.Engine) {
	duration := 5 * time.Second // 减少测试时间
	maxConcurrency := 30        // 减少并发数
	
	fmt.Printf("运行 %v 持续压力测试，最大并发 %d\n", duration, maxConcurrency)
	
	var requestCount int64
	var errorCount int64
	
	stopChan := make(chan struct{})
	time.AfterFunc(duration, func() {
		close(stopChan)
	})
	
	start := time.Now()
	var wg sync.WaitGroup
	
	for i := 0; i < maxConcurrency; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			
			requestID := 0
			for {
				select {
				case <-stopChan:
					return
				default:
					req := createLicenseRequest("/mobaxterm/generate", "POST", 
						workerID*10000+requestID)
					w := httptest.NewRecorder()
					router.ServeHTTP(w, req)
					
					requestCount++
					if w.Code != http.StatusOK && w.Code != http.StatusTooManyRequests {
						errorCount++
					}
					requestID++
				}
			}
		}(i)
	}
	
	wg.Wait()
	elapsed := time.Since(start)
	
	fmt.Printf("总请求数: %d\n", requestCount)
	fmt.Printf("错误数: %d\n", errorCount)
	fmt.Printf("平均QPS: %.2f\n", float64(requestCount)/elapsed.Seconds())
	if requestCount > 0 {
		fmt.Printf("错误率: %.2f%%\n", float64(errorCount)/float64(requestCount)*100)
	}
	
	// 显示最终缓存统计
	showCacheStats(router)
	
	// 显示健康检查
	fmt.Println("\n健康检查:")
	req, _ := http.NewRequest("GET", "/mobaxterm/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	fmt.Printf("健康状态: %s\n", w.Body.String())
}