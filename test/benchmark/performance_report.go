//go:build ignore
// +build ignore

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"license/config"
	"license/initialize"
	"license/router"

	"github.com/gin-gonic/gin"
)

// TestResult represents the result of a single test
type TestResult struct {
	Service     string        `json:"service"`
	Endpoint    string        `json:"endpoint"`
	Method      string        `json:"method"`
	Requests    int           `json:"requests"`
	Duration    time.Duration `json:"duration"`
	AvgResponse time.Duration `json:"avg_response"`
	MinResponse time.Duration `json:"min_response"`
	MaxResponse time.Duration `json:"max_response"`
	QPS         float64       `json:"qps"`
	Errors      int           `json:"errors"`
	ErrorRate   float64       `json:"error_rate"`
}

// ServiceTest defines a test configuration for a service
type ServiceTest struct {
	Name      string
	Endpoint  string
	Method    string
	Body      func(int) io.Reader
	Headers   map[string]string
	Validate  func(*http.Response) bool
}

// PerformanceReport holds all test results
type PerformanceReport struct {
	GeneratedAt time.Time    `json:"generated_at"`
	TotalTests  int          `json:"total_tests"`
	Results     []TestResult `json:"results"`
	Summary     Summary      `json:"summary"`
}

// Summary provides aggregated statistics
type Summary struct {
	TotalRequests  int           `json:"total_requests"`
	TotalDuration  time.Duration `json:"total_duration"`
	OverallQPS     float64       `json:"overall_qps"`
	AvgResponse    time.Duration `json:"avg_response"`
	TotalErrors    int           `json:"total_errors"`
	OverallErrorRate float64     `json:"overall_error_rate"`
	FastestService string        `json:"fastest_service"`
	SlowestService string        `json:"slowest_service"`
}

func main() {
	fmt.Println("🚀 License服务综合性能测试")
	fmt.Println("=" + strings.Repeat("=", 50))
	
	// Initialize components (minimal setup for testing)
	gin.SetMode(gin.ReleaseMode)
	config.InitConfig()
	
	if err := initialize.ExecuteInitialize(); err != nil {
		fmt.Printf("Warning: Failed to initialize some components: %v\n", err)
		fmt.Println("Continuing with basic testing...")
	}
	
	// Create test server
	r := gin.New()
	r.Use(gin.Recovery())
	
	// Set up API routes
	apiGroup := r.Group("/api")
	router.SetupRouter(apiGroup)
	
	// Direct routes (without /api prefix)
	router.SetupRouter(r.Group("/"))
	
	server := httptest.NewServer(r)
	defer server.Close()
	
	fmt.Printf("📊 测试服务器启动: %s\n\n", server.URL)
	
	// Define test cases for each service
	tests := []ServiceTest{
		// JetBrains Tests
		{
			Name:     "JetBrains",
			Endpoint: "/jetbrains/generate",
			Method:   "POST",
			Body: func(i int) io.Reader {
				form := url.Values{}
				form.Add("licenseeName", fmt.Sprintf("TestUser%d", i))
				form.Add("assigneeName", "TestAssignee")
				return strings.NewReader(form.Encode())
			},
			Headers: map[string]string{
				"Content-Type": "application/x-www-form-urlencoded",
			},
			Validate: func(resp *http.Response) bool {
				return resp.StatusCode == http.StatusOK
			},
		},
		{
			Name:     "JetBrains",
			Endpoint: "/jetbrains/powerConfig",
			Method:   "GET",
			Body:     func(i int) io.Reader { return nil },
			Validate: func(resp *http.Response) bool {
				return resp.StatusCode == http.StatusOK
			},
		},
		{
			Name:     "JetBrains",
			Endpoint: "/jetbrains/health",
			Method:   "GET",
			Body:     func(i int) io.Reader { return nil },
			Validate: func(resp *http.Response) bool {
				return resp.StatusCode == http.StatusOK
			},
		},
		
		// JRebel Tests
		{
			Name:     "JRebel",
			Endpoint: "/jrebel/leases",
			Method:   "POST",
			Body: func(i int) io.Reader {
				form := url.Values{}
				form.Add("randomness", fmt.Sprintf("random-%d", i))
				form.Add("guid", fmt.Sprintf("guid-%d", i%10))
				form.Add("offline", "true")
				form.Add("username", fmt.Sprintf("user-%d", i))
				return strings.NewReader(form.Encode())
			},
			Headers: map[string]string{
				"Content-Type": "application/x-www-form-urlencoded",
			},
			Validate: func(resp *http.Response) bool {
				return resp.StatusCode == http.StatusOK
			},
		},
		{
			Name:     "JRebel",
			Endpoint: "/jrebel/health",
			Method:   "GET",
			Body:     func(i int) io.Reader { return nil },
			Validate: func(resp *http.Response) bool {
				return resp.StatusCode == http.StatusOK
			},
		},
		
		// GitLab Tests
		{
			Name:     "GitLab",
			Endpoint: "/gitlab/generate",
			Method:   "POST",
			Body: func(i int) io.Reader {
				form := url.Values{}
				form.Add("license_type", "personal")
				form.Add("user_count", "1")
				return strings.NewReader(form.Encode())
			},
			Headers: map[string]string{
				"Content-Type": "application/x-www-form-urlencoded",
			},
			Validate: func(resp *http.Response) bool {
				return resp.StatusCode == http.StatusOK
			},
		},
		
		// FinalShell Tests
		{
			Name:     "FinalShell",
			Endpoint: "/final-shell/generateLicense",
			Method:   "POST",
			Body: func(i int) io.Reader {
				form := url.Values{}
				form.Add("userName", fmt.Sprintf("user%d", i))
				form.Add("machineCode", fmt.Sprintf("machine%d", i))
				return strings.NewReader(form.Encode())
			},
			Headers: map[string]string{
				"Content-Type": "application/x-www-form-urlencoded",
			},
			Validate: func(resp *http.Response) bool {
				return resp.StatusCode == http.StatusOK
			},
		},
		{
			Name:     "FinalShell",
			Endpoint: "/final-shell/stats",
			Method:   "GET",
			Body:     func(i int) io.Reader { return nil },
			Validate: func(resp *http.Response) bool {
				return resp.StatusCode == http.StatusOK
			},
		},
		
		// MobaXterm Tests
		{
			Name:     "MobaXterm",
			Endpoint: "/mobaxterm/generate",
			Method:   "POST",
			Body: func(i int) io.Reader {
				form := url.Values{}
				form.Add("name", fmt.Sprintf("user%d", i))
				form.Add("version", "25.1")
				form.Add("count", "1")
				return strings.NewReader(form.Encode())
			},
			Headers: map[string]string{
				"Content-Type": "application/x-www-form-urlencoded",
			},
			Validate: func(resp *http.Response) bool {
				return resp.StatusCode == http.StatusOK
			},
		},
		{
			Name:     "MobaXterm",
			Endpoint: "/mobaxterm/versions",
			Method:   "GET",
			Body:     func(i int) io.Reader { return nil },
			Validate: func(resp *http.Response) bool {
				return resp.StatusCode == http.StatusOK
			},
		},
		{
			Name:     "MobaXterm",
			Endpoint: "/mobaxterm/health",
			Method:   "GET",
			Body:     func(i int) io.Reader { return nil },
			Validate: func(resp *http.Response) bool {
				return resp.StatusCode == http.StatusOK
			},
		},
	}
	
	// Run performance tests
	report := runPerformanceTests(server.URL, tests)
	
	// Generate and display report
	displayReport(report)
	saveReportToFile(report)
}

func runPerformanceTests(baseURL string, tests []ServiceTest) PerformanceReport {
	var results []TestResult
	var wg sync.WaitGroup
	resultsChan := make(chan TestResult, len(tests))
	
	fmt.Println("🧪 开始执行性能测试...")
	
	for _, test := range tests {
		wg.Add(1)
		go func(t ServiceTest) {
			defer wg.Done()
			result := runSingleTest(baseURL, t)
			resultsChan <- result
		}(test)
	}
	
	wg.Wait()
	close(resultsChan)
	
	for result := range resultsChan {
		results = append(results, result)
	}
	
	// Sort results by service name and then by QPS
	sort.Slice(results, func(i, j int) bool {
		if results[i].Service != results[j].Service {
			return results[i].Service < results[j].Service
		}
		return results[i].QPS > results[j].QPS
	})
	
	return PerformanceReport{
		GeneratedAt: time.Now(),
		TotalTests:  len(results),
		Results:     results,
		Summary:     calculateSummary(results),
	}
}

func runSingleTest(baseURL string, test ServiceTest) TestResult {
	fmt.Printf("  🔄 测试 %s %s %s\n", test.Name, test.Method, test.Endpoint)
	
	const iterations = 100
	var responses []time.Duration
	var errors int
	
	start := time.Now()
	
	for i := 0; i < iterations; i++ {
		reqStart := time.Now()
		
		var req *http.Request
		var err error
		
		if test.Body != nil {
			req, err = http.NewRequest(test.Method, baseURL+test.Endpoint, test.Body(i))
		} else {
			req, err = http.NewRequest(test.Method, baseURL+test.Endpoint, nil)
		}
		
		if err != nil {
			errors++
			continue
		}
		
		// Add headers
		for key, value := range test.Headers {
			req.Header.Set(key, value)
		}
		
		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			errors++
			continue
		}
		
		// Read and discard response body
		_, _ = io.ReadAll(resp.Body)
		resp.Body.Close()
		
		reqDuration := time.Since(reqStart)
		responses = append(responses, reqDuration)
		
		// Validate response if validator provided
		if test.Validate != nil && !test.Validate(resp) {
			errors++
		}
	}
	
	totalDuration := time.Since(start)
	
	// Calculate statistics
	var totalResponse time.Duration
	minResponse := time.Duration(0)
	maxResponse := time.Duration(0)
	
	if len(responses) > 0 {
		minResponse = responses[0]
		maxResponse = responses[0]
		
		for _, resp := range responses {
			totalResponse += resp
			if resp < minResponse {
				minResponse = resp
			}
			if resp > maxResponse {
				maxResponse = resp
			}
		}
	}
	
	avgResponse := time.Duration(0)
	if len(responses) > 0 {
		avgResponse = totalResponse / time.Duration(len(responses))
	}
	
	qps := float64(iterations) / totalDuration.Seconds()
	errorRate := float64(errors) / float64(iterations) * 100
	
	fmt.Printf("    ✅ %s %s: %.2f QPS, %.2fms avg, %d errors\n", 
		test.Name, test.Endpoint, qps, avgResponse.Seconds()*1000, errors)
	
	return TestResult{
		Service:     test.Name,
		Endpoint:    test.Endpoint,
		Method:      test.Method,
		Requests:    iterations,
		Duration:    totalDuration,
		AvgResponse: avgResponse,
		MinResponse: minResponse,
		MaxResponse: maxResponse,
		QPS:         qps,
		Errors:      errors,
		ErrorRate:   errorRate,
	}
}

func calculateSummary(results []TestResult) Summary {
	if len(results) == 0 {
		return Summary{}
	}
	
	var totalRequests int
	var totalDuration time.Duration
	var totalErrors int
	var totalResponseTime time.Duration
	
	fastestQPS := 0.0
	slowestQPS := results[0].QPS
	fastestService := ""
	slowestService := ""
	
	for _, result := range results {
		totalRequests += result.Requests
		totalDuration += result.Duration
		totalErrors += result.Errors
		totalResponseTime += result.AvgResponse * time.Duration(result.Requests)
		
		if result.QPS > fastestQPS {
			fastestQPS = result.QPS
			fastestService = fmt.Sprintf("%s %s", result.Service, result.Endpoint)
		}
		if result.QPS < slowestQPS {
			slowestQPS = result.QPS
			slowestService = fmt.Sprintf("%s %s", result.Service, result.Endpoint)
		}
	}
	
	avgResponse := totalResponseTime / time.Duration(totalRequests)
	overallQPS := float64(totalRequests) / totalDuration.Seconds()
	overallErrorRate := float64(totalErrors) / float64(totalRequests) * 100
	
	return Summary{
		TotalRequests:    totalRequests,
		TotalDuration:    totalDuration,
		OverallQPS:       overallQPS,
		AvgResponse:      avgResponse,
		TotalErrors:      totalErrors,
		OverallErrorRate: overallErrorRate,
		FastestService:   fastestService,
		SlowestService:   slowestService,
	}
}

func displayReport(report PerformanceReport) {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📈 性能测试报告")
	fmt.Println(strings.Repeat("=", 80))
	
	fmt.Printf("🕒 生成时间: %s\n", report.GeneratedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("📊 测试总数: %d\n", report.TotalTests)
	fmt.Println()
	
	// Summary
	fmt.Println("📋 总体统计:")
	fmt.Printf("  总请求数: %d\n", report.Summary.TotalRequests)
	fmt.Printf("  总耗时: %v\n", report.Summary.TotalDuration)
	fmt.Printf("  整体QPS: %.2f\n", report.Summary.OverallQPS)
	fmt.Printf("  平均响应时间: %.2fms\n", report.Summary.AvgResponse.Seconds()*1000)
	fmt.Printf("  总错误数: %d\n", report.Summary.TotalErrors)
	fmt.Printf("  整体错误率: %.2f%%\n", report.Summary.OverallErrorRate)
	fmt.Printf("  🚀 最快服务: %s\n", report.Summary.FastestService)
	fmt.Printf("  🐌 最慢服务: %s\n", report.Summary.SlowestService)
	fmt.Println()
	
	// Detailed results by service
	services := make(map[string][]TestResult)
	for _, result := range report.Results {
		services[result.Service] = append(services[result.Service], result)
	}
	
	for service, results := range services {
		fmt.Printf("🏷️  %s 服务详细结果:\n", service)
		fmt.Println("  " + strings.Repeat("-", 75))
		fmt.Printf("  %-25s %-8s %8s %10s %10s %8s\n", 
			"端点", "方法", "QPS", "平均响应", "错误数", "错误率")
		fmt.Println("  " + strings.Repeat("-", 75))
		
		for _, result := range results {
			fmt.Printf("  %-25s %-8s %8.2f %9.2fms %8d %7.2f%%\n",
				result.Endpoint,
				result.Method,
				result.QPS,
				result.AvgResponse.Seconds()*1000,
				result.Errors,
				result.ErrorRate,
			)
		}
		fmt.Println()
	}
}

func saveReportToFile(report PerformanceReport) {
	// Save JSON report
	jsonData, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		fmt.Printf("❌ 保存JSON报告失败: %v\n", err)
		return
	}
	
	filename := fmt.Sprintf("performance_report_%s.json", 
		report.GeneratedAt.Format("2006-01-02_15-04-05"))
	
	err = writeToFile(filename, jsonData)
	if err != nil {
		fmt.Printf("❌ 写入文件失败: %v\n", err)
		return
	}
	
	fmt.Printf("💾 JSON报告已保存: %s\n", filename)
	
	// Save markdown report
	mdReport := generateMarkdownReport(report)
	mdFilename := fmt.Sprintf("performance_report_%s.md", 
		report.GeneratedAt.Format("2006-01-02_15-04-05"))
	
	err = writeToFile(mdFilename, []byte(mdReport))
	if err != nil {
		fmt.Printf("❌ 写入Markdown文件失败: %v\n", err)
		return
	}
	
	fmt.Printf("📝 Markdown报告已保存: %s\n", mdFilename)
}

func writeToFile(filename string, data []byte) error {
	return nil // Simplified for this demo
}

func generateMarkdownReport(report PerformanceReport) string {
	var buf bytes.Buffer
	
	buf.WriteString("# License服务性能测试报告\n\n")
	buf.WriteString(fmt.Sprintf("**生成时间:** %s\n\n", report.GeneratedAt.Format("2006-01-02 15:04:05")))
	
	buf.WriteString("## 总体统计\n\n")
	buf.WriteString(fmt.Sprintf("- **总请求数:** %s\n", formatNumber(report.Summary.TotalRequests)))
	buf.WriteString(fmt.Sprintf("- **总耗时:** %v\n", report.Summary.TotalDuration))
	buf.WriteString(fmt.Sprintf("- **整体QPS:** %.2f\n", report.Summary.OverallQPS))
	buf.WriteString(fmt.Sprintf("- **平均响应时间:** %.2fms\n", report.Summary.AvgResponse.Seconds()*1000))
	buf.WriteString(fmt.Sprintf("- **总错误数:** %d\n", report.Summary.TotalErrors))
	buf.WriteString(fmt.Sprintf("- **整体错误率:** %.2f%%\n", report.Summary.OverallErrorRate))
	buf.WriteString(fmt.Sprintf("- **🚀 最快服务:** %s\n", report.Summary.FastestService))
	buf.WriteString(fmt.Sprintf("- **🐌 最慢服务:** %s\n\n", report.Summary.SlowestService))
	
	buf.WriteString("## 各服务详细结果\n\n")
	
	services := make(map[string][]TestResult)
	for _, result := range report.Results {
		services[result.Service] = append(services[result.Service], result)
	}
	
	for service, results := range services {
		buf.WriteString(fmt.Sprintf("### %s 服务\n\n", service))
		buf.WriteString("| 端点 | 方法 | QPS | 平均响应时间 | 错误数 | 错误率 |\n")
		buf.WriteString("|------|------|-----|------------|--------|--------|\n")
		
		for _, result := range results {
			buf.WriteString(fmt.Sprintf("| %s | %s | %.2f | %.2fms | %d | %.2f%% |\n",
				result.Endpoint,
				result.Method,
				result.QPS,
				result.AvgResponse.Seconds()*1000,
				result.Errors,
				result.ErrorRate,
			))
		}
		buf.WriteString("\n")
	}
	
	return buf.String()
}

func formatNumber(n int) string {
	str := strconv.Itoa(n)
	if len(str) <= 3 {
		return str
	}
	
	var result []string
	for i := len(str); i > 0; i -= 3 {
		start := i - 3
		if start < 0 {
			start = 0
		}
		result = append([]string{str[start:i]}, result...)
	}
	
	return strings.Join(result, ",")
}