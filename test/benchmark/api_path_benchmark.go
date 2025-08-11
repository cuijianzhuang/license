//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"strings"
	"time"

	"license/config"
	"license/initialize"
	"license/router"

	"github.com/gin-gonic/gin"
)

// 模拟旧的API路径检测方法用于对比
func oldIsAPIPath(path string) bool {
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
	fmt.Println("🔄 API路径检测性能对比测试")
	fmt.Println("=" + strings.Repeat("=", 60))
	
	gin.SetMode(gin.ReleaseMode)
	config.InitConfig()
	
	if err := initialize.ExecuteInitialize(); err != nil {
		fmt.Printf("Warning: %v\n", err)
	}
	
	testPaths := []string{
		// API 路径
		"/jetbrains/generate",
		"/jetbrains/health",
		"/jetbrains/powerConfig",
		"/jrebel/leases",
		"/jrebel/health",
		"/mobaxterm/generate",
		"/mobaxterm/versions",
		"/final-shell/generateLicense",
		"/final-shell/stats",
		"/gitlab/generate",
		"/server/status",
		"/agent/leases",
		// 非API路径
		"/static/js/main.js",
		"/static/css/app.css",
		"/index.html",
		"/favicon.ico",
		"/dashboard",
		"/settings",
		"/home",
		"/",
	}
	
	iterations := 1000000
	
	fmt.Printf("📊 测试路径数: %d\n", len(testPaths))
	fmt.Printf("🔄 测试迭代数: %s\n\n", formatNumber(iterations))
	
	// 测试旧方法
	fmt.Println("🕰️  测试旧的API路径检测方法...")
	start := time.Now()
	for i := 0; i < iterations; i++ {
		for _, path := range testPaths {
			oldIsAPIPath(path)
		}
	}
	oldDuration := time.Since(start)
	oldAvg := oldDuration / time.Duration(iterations*len(testPaths))
	
	// 测试新方法
	fmt.Println("⚡ 测试优化后的API路径检测方法...")
	start = time.Now()
	for i := 0; i < iterations; i++ {
		for _, path := range testPaths {
			router.IsAPIPath(path)
		}
	}
	newDuration := time.Since(start)
	newAvg := newDuration / time.Duration(iterations*len(testPaths))
	
	// 计算性能提升
	improvement := float64(oldDuration) / float64(newDuration)
	
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📈 API路径检测性能测试结果")
	fmt.Println(strings.Repeat("=", 80))
	
	fmt.Printf("🕰️  旧方法 (String Slice遍历):\n")
	fmt.Printf("   总耗时: %v\n", oldDuration)
	fmt.Printf("   平均耗时: %v 每次操作\n", oldAvg)
	fmt.Printf("   QPS: %.0f\n\n", float64(iterations*len(testPaths))/oldDuration.Seconds())
	
	fmt.Printf("⚡ 新方法 (Map O(1)查找):\n")
	fmt.Printf("   总耗时: %v\n", newDuration)
	fmt.Printf("   平均耗时: %v 每次操作\n", newAvg)
	fmt.Printf("   QPS: %.0f\n\n", float64(iterations*len(testPaths))/newDuration.Seconds())
	
	fmt.Printf("🚀 性能提升: %.2fx 更快\n", improvement)
	fmt.Printf("⏱️  时间节省: %v (%.1f%%)\n", 
		oldDuration-newDuration, 
		float64(oldDuration-newDuration)/float64(oldDuration)*100)
	
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("✅ 正确性验证")
	fmt.Println(strings.Repeat("=", 80))
	
	allCorrect := true
	for _, path := range testPaths {
		oldResult := oldIsAPIPath(path)
		newResult := router.IsAPIPath(path)
		
		status := "✅"
		if oldResult != newResult {
			status = "❌"
			allCorrect = false
		}
		
		fmt.Printf("%s %-30s -> %-5v (旧: %-5v, 新: %-5v)\n", 
			status, path, oldResult, oldResult, newResult)
	}
	
	fmt.Println(strings.Repeat("-", 80))
	if allCorrect {
		fmt.Println("🎉 所有测试通过！优化后的方法与原方法结果完全一致。")
	} else {
		fmt.Println("⚠️  警告：优化后的方法与原方法结果不一致！")
	}
	
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📊 综合性能报告")
	fmt.Println(strings.Repeat("=", 80))
	
	fmt.Printf("🎯 优化目标: API路径检测性能提升\n")
	fmt.Printf("📈 实际提升: %.2fx 性能改善\n", improvement)
	fmt.Printf("💡 优化方式: String Slice遍历 → Map O(1)查找\n")
	fmt.Printf("🔢 测试规模: %s 次路径检测操作\n", formatNumber(iterations*len(testPaths)))
	fmt.Printf("⏱️  节省时间: 每百万次操作节省 %.2fms\n", 
		float64(oldDuration-newDuration)/float64(iterations)*1000000/1000000)
	fmt.Printf("📦 内存开销: Map存储 vs Slice遍历 (略有增加，但换来显著性能提升)\n")
	fmt.Printf("✨ 建议: 该优化在高并发场景下效果更明显\n")
}

func formatNumber(n int) string {
	s := fmt.Sprintf("%d", n)
	if len(s) <= 3 {
		return s
	}
	
	var result []string
	for i := len(s); i > 0; i -= 3 {
		start := i - 3
		if start < 0 {
			start = 0
		}
		result = append([]string{s[start:i]}, result...)
	}
	
	return strings.Join(result, ",")
}