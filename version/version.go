package version

// 默认版本，可以在编译时通过 -ldflags="-X 'license/version.Version=x.y.z'" 注入
var (
	// Version 存储当前应用程序版本
	Version = "0.0.1"
	
	// BuildTime 存储构建时间
	BuildTime = "unknown"
	
	// GitCommit 存储Git提交哈希值
	GitCommit = "unknown"
)

// GetVersion 返回应用程序版本
func GetVersion() string {
	return Version
}

// GetBuildInfo 返回包含版本、构建时间和Git提交信息的映射
func GetBuildInfo() map[string]string {
	return map[string]string{
		"version":   Version,
		"buildTime": BuildTime,
		"gitCommit": GitCommit,
	}
} 