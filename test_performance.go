package main

import (
	"bytes"
	"fmt"
	"license/gitlab/entity"
	"license/gitlab/service"
	"net/http"
	"net/http/httptest"
	"net/url"
	"runtime"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

// Performance test for GitLab license generation
func main() {
	fmt.Println("GitLab License Generation Performance Test")
	fmt.Println("=========================================")

	// Initialize the service
	if err := service.LoadKeys(); err != nil {
		fmt.Printf("Warning: Failed to load keys (this is expected in test environment): %v\n", err)
	}

	// Set Gin to release mode for accurate performance testing
	gin.SetMode(gin.ReleaseMode)

	// Create a test router
	r := gin.New()
	r.POST("/gitlab/generate", func(c *gin.Context) {
		licenseInfo := entity.LicenseInfo{
			Name:    "Test User",
			Email:   "test@example.com",
			Company: "Test Company",
		}
		service.Generate(c, licenseInfo, "")
	})

	// Run performance tests
	runPerformanceTest(r)
}

func runPerformanceTest(router *gin.Engine) {
	numRequests := 100
	concurrency := 10

	fmt.Printf("Running %d requests with concurrency %d...\n", numRequests, concurrency)

	// Get initial memory stats
	var m1, m2 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m1)

	start := time.Now()

	// Create channels for controlling concurrency
	semaphore := make(chan struct{}, concurrency)
	results := make(chan time.Duration, numRequests)

	// Launch goroutines
	for i := 0; i < numRequests; i++ {
		go func() {
			semaphore <- struct{}{}        // Acquire semaphore
			defer func() { <-semaphore }() // Release semaphore

			reqStart := time.Now()

			// Create test request
			form := url.Values{}
			form.Add("Name", "Test User")
			form.Add("Email", "test@example.com")
			form.Add("Company", "Test Company")

			req, _ := http.NewRequest("POST", "/gitlab/generate", bytes.NewBufferString(form.Encode()))
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			results <- time.Since(reqStart)
		}()
	}

	// Collect results
	var totalTime time.Duration
	var minTime, maxTime time.Duration
	cacheHits := 0

	for i := 0; i < numRequests; i++ {
		reqTime := <-results
		totalTime += reqTime

		if i == 0 {
			minTime = reqTime
			maxTime = reqTime
		} else {
			if reqTime < minTime {
				minTime = reqTime
			}
			if reqTime > maxTime {
				maxTime = reqTime
			}
		}
	}

	elapsed := time.Since(start)

	// Get final memory stats
	runtime.GC()
	runtime.ReadMemStats(&m2)

	// Print results
	fmt.Printf("\nPerformance Results:\n")
	fmt.Printf("Total time: %v\n", elapsed)
	fmt.Printf("Average response time: %v\n", totalTime/time.Duration(numRequests))
	fmt.Printf("Min response time: %v\n", minTime)
	fmt.Printf("Max response time: %v\n", maxTime)
	fmt.Printf("Requests per second: %.2f\n", float64(numRequests)/elapsed.Seconds())

	fmt.Printf("\nMemory Usage:\n")
	fmt.Printf("Memory allocated during test: %d KB\n", (m2.TotalAlloc-m1.TotalAlloc)/1024)
	fmt.Printf("Memory in use after test: %d KB\n", m2.Alloc/1024)
	fmt.Printf("GC runs during test: %d\n", m2.NumGC-m1.NumGC)

	// Test cache effectiveness
	fmt.Printf("\nTesting Cache Effectiveness...\n")
	testCacheEffectiveness(router)
}

func testCacheEffectiveness(router *gin.Engine) {
	// Make the same request twice to test caching
	form := url.Values{}
	form.Add("Name", "Cache Test User")
	form.Add("Email", "cache@example.com")
	form.Add("Company", "Cache Test Company")

	// First request (cache miss)
	start := time.Now()
	req1, _ := http.NewRequest("POST", "/gitlab/generate", bytes.NewBufferString(form.Encode()))
	req1.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)
	firstReqTime := time.Since(start)

	// Second request (cache hit)
	start = time.Now()
	req2, _ := http.NewRequest("POST", "/gitlab/generate", bytes.NewBufferString(form.Encode()))
	req2.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)
	secondReqTime := time.Since(start)

	fmt.Printf("First request (cache miss): %v\n", firstReqTime)
	fmt.Printf("Second request (cache hit): %v\n", secondReqTime)
	if secondReqTime > 0 && firstReqTime > 0 {
		speedup := float64(firstReqTime) / float64(secondReqTime)
		fmt.Printf("Cache speedup: %.2fx faster\n", speedup)
	}

	// Check cache headers
	cacheHeader1 := w1.Header().Get("X-Cache")
	cacheHeader2 := w2.Header().Get("X-Cache")
	fmt.Printf("First request cache header: %s\n", cacheHeader1)
	fmt.Printf("Second request cache header: %s\n", cacheHeader2)
}
