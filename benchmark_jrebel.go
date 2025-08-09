// +build ignore

package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"license/jrebel/api"
	"license/jrebel/constant"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("JRebel 许可证服务性能基准测试")
	fmt.Println("================================")

	gin.SetMode(gin.ReleaseMode)

	// 创建原始控制器
	originalController := api.NewLeasesController()

	// 创建优化控制器
	optimizedController, err := api.NewOptimizedLeasesController()
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

func setupOriginalRouter(controller *api.LeasesController) *gin.Engine {
	r := gin.New()
	r.POST("/jrebel/leases", controller.LeasesHandler)
	r.DELETE("/jrebel/leases/1", controller.Leases1Handler)
	r.POST("/jrebel/validate-connection", controller.ValidateHandler)
	return r
}

func setupOptimizedRouter(controller *api.OptimizedLeasesController) *gin.Engine {
	r := gin.New()
	r.POST("/jrebel/leases", controller.OptimizedLeasesHandler)
	r.DELETE("/jrebel/leases/1", controller.OptimizedLeases1Handler)
	r.POST("/jrebel/validate-connection", controller.OptimizedValidateHandler)
	r.GET("/jrebel/stats", controller.GetPerformanceStats)
	r.POST("/jrebel/clear-cache", controller.ClearCache)
	return r
}

func runBenchmarkTests(originalRouter, optimizedRouter *gin.Engine) {
	testCases := []struct {
		name     string
		endpoint string
		method   string
	}{
		{"Leases Handler", "/jrebel/leases", "POST"},
		{"Leases/1 Handler", "/jrebel/leases/1", "DELETE"},
		{"Validate Handler", "/jrebel/validate-connection", "POST"},
	}

	for _, tc := range testCases {
		fmt.Printf("\n=== %s 性能测试 ===\n", tc.name)
		
		// 原始版本测试
		fmt.Println("原始版本:")
		originalTime := benchmarkEndpoint(originalRouter, tc.endpoint, tc.method, 100)
		
		// 优化版本测试（无缓存）
		clearCache(optimizedRouter)
		fmt.Println("优化版本（无缓存）:")
		optimizedNoCacheTime := benchmarkEndpoint(optimizedRouter, tc.endpoint, tc.method, 100)
		
		// 优化版本测试（有缓存）
		preloadCache(optimizedRouter, tc.endpoint, tc.method, 10)
		fmt.Println("优化版本（有缓存）:")
		optimizedCachedTime := benchmarkEndpoint(optimizedRouter, tc.endpoint, tc.method, 100)
		
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

func benchmarkEndpoint(router *gin.Engine, endpoint, method string, iterations int) time.Duration {
	start := time.Now()
	
	for i := 0; i < iterations; i++ {
		req := createTestRequest(endpoint, method, i)
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

func createTestRequest(endpoint, method string, seed int) *http.Request {
	form := url.Values{}
	form.Add("randomness", fmt.Sprintf("random-%d", seed))
	form.Add("guid", fmt.Sprintf("guid-%d", seed%10)) // 10种不同的GUID循环
	form.Add("offline", "true")
	form.Add("clientTime", strconv.FormatInt(time.Now().UnixMilli(), 10))
	form.Add("username", fmt.Sprintf("user-%d", seed))

	body := strings.NewReader(form.Encode())
	req, _ := http.NewRequest(method, endpoint, body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	
	return req
}

func preloadCache(router *gin.Engine, endpoint, method string, count int) {
	for i := 0; i < count; i++ {
		req := createTestRequest(endpoint, method, i)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}

func clearCache(router *gin.Engine) {
	req, _ := http.NewRequest("POST", "/jrebel/clear-cache", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
}

func showCacheStats(router *gin.Engine) {
	req, _ := http.NewRequest("GET", "/jrebel/stats", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	fmt.Printf("缓存统计: %s\n", w.Body.String())
}

func runConcurrencyTest(originalRouter, optimizedRouter *gin.Engine) {
	concurrencyLevels := []int{10, 50, 100}
	requestsPerLevel := 200
	
	for _, concurrency := range concurrencyLevels {
		fmt.Printf("\n并发数: %d\n", concurrency)
		
		// 原始版本
		originalTime := runConcurrentRequests(originalRouter, "/jrebel/leases", "POST", 
			concurrency, requestsPerLevel)
		
		// 优化版本
		clearCache(optimizedRouter)
		optimizedTime := runConcurrentRequests(optimizedRouter, "/jrebel/leases", "POST", 
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
				req := createTestRequest(endpoint, method, workerID*requestsPerWorker+j)
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)
			}
		}(i)
	}
	
	wg.Wait()
	return time.Since(start)
}

func runStressTest(optimizedRouter *gin.Engine) {
	duration := 10 * time.Second
	maxConcurrency := 200
	
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
					req := createTestRequest("/jrebel/leases", "POST", 
						workerID*10000+requestID)
					w := httptest.NewRecorder()
					optimizedRouter.ServeHTTP(w, req)
					
					requestCount++
					if w.Code != http.StatusOK {
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
}

// 原始签名函数用于对比测试
func originalSign(clientRandomness, guid string, offline bool, validFrom, validUntil int64) string {
	signatureBase := clientRandomness + ";" + constant.ServerRandomness + ";" + guid + ";" + strconv.FormatBool(offline)
	if offline {
		signatureBase += ";" + strconv.FormatInt(validFrom, 10) + ";" + strconv.FormatInt(validUntil, 10)
	}

	block, _ := pem.Decode([]byte(constant.LeasesPrivateKey))
	if block == nil {
		return ""
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return ""
	}

	hash := sha1.New()
	hash.Write([]byte(signatureBase))
	hashed := hash.Sum(nil)

	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA1, hashed)
	if err != nil {
		return ""
	}

	return base64.StdEncoding.EncodeToString(signature)
}