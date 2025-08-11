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

	// 创建原始控制器
	originalController := api.NewMobaXtermController()

	// 创建优化控制器
	optimizedController, err := api.NewOptimizedController()
	if err != nil {
		fmt.Printf("Failed to create optimized controller: %v\n", err)
		return
	}
	defer optimizedController.Shutdown()

	// 创建路由
	originalRouter := setupOriginalRouter(originalController)
	optimizedRouter := setupOptimizedRouter(optimizedController)

	// 运行基准测试
	runBenchmarkTests(originalRouter, optimizedRouter)
}

func setupOriginalRouter(controller *api.Controller) *gin.Engine {
	r := gin.New()
	r.GET("/mobaxterm/versions", controller.FetchVersions)
	r.POST("/mobaxterm/generate", controller.GenerateLicense)
	r.GET("/mobaxterm/generate", controller.GenerateLicense)
	return r
}

func setupOptimizedRouter(controller *api.OptimizedController) *gin.Engine {
	r := gin.New()
	// 应用限流中间件
	r.Use(controller.RateLimitMiddleware())

	r.GET("/mobaxterm/versions", controller.OptimizedFetchVersions)
	r.POST("/mobaxterm/generate", controller.OptimizedGenerateLicense)
	r.GET("/mobaxterm/generate", controller.OptimizedGenerateLicense)
	r.GET("/mobaxterm/stats", controller.GetPerformanceStats)
	r.POST("/mobaxterm/clear-cache", controller.ClearCache)
	r.GET("/mobaxterm/health", controller.HealthCheck)
	r.POST("/mobaxterm/force-gc", controller.ForceGC)
	return r
}

func runBenchmarkTests(originalRouter, optimizedRouter *gin.Engine) {
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

		// 原始版本测试
		fmt.Println("原始版本:")
		originalTime := tc.testFunc(originalRouter, tc.endpoint, tc.method, 100)

		// 优化版本测试（无缓存）
		if tc.name == "License Generation" {
			clearCache(optimizedRouter)
		}
		fmt.Println("优化版本（无缓存）:")
		optimizedNoCacheTime := tc.testFunc(optimizedRouter, tc.endpoint, tc.method, 100)

		// 优化版本测试（有缓存）
		if tc.name == "License Generation" {
			preloadLicenseCache(optimizedRouter, 10)
		}
		fmt.Println("优化版本（有缓存）:")
		optimizedCachedTime := tc.testFunc(optimizedRouter, tc.endpoint, tc.method, 100)

		// 计算提升
		noCacheImprovement := float64(originalTime) / float64(optimizedNoCacheTime)
		cachedImprovement := float64(originalTime) / float64(optimizedCachedTime)

		fmt.Printf("无缓存优化倍数: %.2fx\n", noCacheImprovement)
		fmt.Printf("缓存优化倍数: %.2fx\n", cachedImprovement)

		// 显示缓存统计
		showCacheStats(optimizedRouter)
	}

	// 并发测试
	fmt.Println("\n=== 并发性能测试 ===")
	runConcurrencyTest(originalRouter, optimizedRouter)

	// 压力测试
	fmt.Println("\n=== 压力测试 ===")
	runStressTest(optimizedRouter)
}

func benchmarkLicenseGeneration(router *gin.Engine, endpoint, method string, iterations int) time.Duration {
	start := time.Now()

	for i := 0; i < iterations; i++ {
		req := createLicenseRequest(endpoint, method, i)
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

func runConcurrencyTest(originalRouter, optimizedRouter *gin.Engine) {
	concurrencyLevels := []int{10, 50, 100}
	requestsPerLevel := 100 // 减少请求数以避免生成太多许可证文件

	for _, concurrency := range concurrencyLevels {
		fmt.Printf("\n并发数: %d\n", concurrency)

		// 原始版本
		originalTime := runConcurrentRequests(originalRouter, "/mobaxterm/generate", "POST",
			concurrency, requestsPerLevel)

		// 优化版本
		clearCache(optimizedRouter)
		optimizedTime := runConcurrentRequests(optimizedRouter, "/mobaxterm/generate", "POST",
			concurrency, requestsPerLevel)

		improvement := float64(originalTime) / float64(optimizedTime)

		fmt.Printf("原始版本: %v (QPS: %.2f)\n", originalTime,
			float64(concurrency*requestsPerLevel)/originalTime.Seconds())
		fmt.Printf("优化版本: %v (QPS: %.2f)\n", optimizedTime,
			float64(concurrency*requestsPerLevel)/optimizedTime.Seconds())
		fmt.Printf("性能提升: %.2fx\n", improvement)
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

func runStressTest(optimizedRouter *gin.Engine) {
	duration := 5 * time.Second // 减少测试时间
	maxConcurrency := 50        // 减少并发数

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
					optimizedRouter.ServeHTTP(w, req)

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
	fmt.Printf("错误率: %.2f%%\n", float64(errorCount)/float64(requestCount)*100)

	// 显示最终缓存统计
	showCacheStats(optimizedRouter)

	// 显示健康检查
	fmt.Println("\n健康检查:")
	req, _ := http.NewRequest("GET", "/mobaxterm/health", nil)
	w := httptest.NewRecorder()
	optimizedRouter.ServeHTTP(w, req)
	fmt.Printf("健康状态: %s\n", w.Body.String())
}
