//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"license/router"
	"strings"
	"time"
)

// Original implementation for comparison
func originalIsAPIPath(path string) bool {
	apiPrefixes := []string{
		"/server/",
		"/final-shell/",
		"/gitlab/",
		"/rpc/",
		"/jrebel/",
		"/agent/",
		"/mobaxterm/",
		"/jetbrains/",
	}
	
	for _, prefix := range apiPrefixes {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}
	return false
}

func main() {
	fmt.Println("Router Performance Benchmark")
	fmt.Println("============================")
	
	// Test paths
	testPaths := []string{
		"/server/status",
		"/jrebel/leases",
		"/mobaxterm/generate",
		"/jetbrains/generate",
		"/final-shell/generateLicense",
		"/gitlab/generate",
		"/rpc/ping.action",
		"/agent/leases",
		"/static/js/main.js",
		"/index.html",
		"/favicon.ico",
		"/dashboard",
		"/settings",
		"/api/server/status",
	}
	
	iterations := 100000
	
	fmt.Printf("Testing %d iterations for each method...\n\n", iterations)
	
	// Test original implementation
	fmt.Println("Original Implementation:")
	start := time.Now()
	for i := 0; i < iterations; i++ {
		for _, path := range testPaths {
			originalIsAPIPath(path)
		}
	}
	originalTime := time.Since(start)
	fmt.Printf("Time: %v (avg: %v per operation)\n", originalTime, originalTime/time.Duration(iterations*len(testPaths)))
	
	// Test optimized implementation
	fmt.Println("\nOptimized Implementation (Trie):")
	start = time.Now()
	for i := 0; i < iterations; i++ {
		for _, path := range testPaths {
			router.IsAPIPath(path)
		}
	}
	optimizedTime := time.Since(start)
	fmt.Printf("Time: %v (avg: %v per operation)\n", optimizedTime, optimizedTime/time.Duration(iterations*len(testPaths)))
	
	// Calculate improvement
	improvement := float64(originalTime) / float64(optimizedTime)
	fmt.Printf("\nPerformance improvement: %.2fx faster\n", improvement)
	
	// Test correctness
	fmt.Println("\nCorrectness Test:")
	allCorrect := true
	for _, path := range testPaths {
		original := originalIsAPIPath(path)
		optimized := router.IsAPIPath(path)
		
		if original != optimized {
			fmt.Printf("MISMATCH for %s: original=%v, optimized=%v\n", path, original, optimized)
			allCorrect = false
		} else {
			fmt.Printf("✓ %s -> %v\n", path, optimized)
		}
	}
	
	if allCorrect {
		fmt.Println("\n✅ All tests passed - optimization is correct!")
	} else {
		fmt.Println("\n❌ Some tests failed - optimization needs fixing")
	}
}