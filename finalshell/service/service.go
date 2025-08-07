package service

import (
	"crypto/md5"
	"hash"
	"sync"
	"time"
	"unsafe"

	"github.com/ebfe/keccak"
)

// 缓存配置
const (
	cacheMaxSize = 1000
	cacheTTL     = 30 * time.Minute
)

// 缓存项
type cacheItem struct {
	license   FinalShellLicense
	timestamp time.Time
}

// 全局缓存和对象池
var (
	licenseCache = make(map[string]*cacheItem)
	cacheMutex   sync.RWMutex
	
	// 预分配的字符串构建器池
	builderPool = sync.Pool{
		New: func() interface{} {
			builder := make([]byte, 0, 64) // 预分配64字节
			return &builder
		},
	}
	
	// 预分配的哈希实例池
	md5Pool = sync.Pool{
		New: func() interface{} {
			return md5.New()
		},
	}
	
	keccakPool = sync.Pool{
		New: func() interface{} {
			return keccak.New384()
		},
	}
)

// 优化的哈希函数使用对象池
func md5HashOptimized(parts ...string) string {
	hasher := md5Pool.Get().(hash.Hash)
	defer func() {
		hasher.Reset()
		md5Pool.Put(hasher)
	}()

	// 使用字节切片池优化字符串拼接
	builder := builderPool.Get().(*[]byte)
	defer func() {
		*builder = (*builder)[:0] // 重置切片长度但保留容量
		builderPool.Put(builder)
	}()

	// 预计算总长度以避免重复分配
	totalLen := 0
	for _, part := range parts {
		totalLen += len(part)
	}
	
	if cap(*builder) < totalLen {
		*builder = make([]byte, 0, totalLen)
	}

	// 高效字符串拼接
	for _, part := range parts {
		*builder = append(*builder, part...)
	}

	hasher.Write(*builder)
	hash := hasher.Sum(nil)
	
	// 直接操作字节避免hex.EncodeToString的开销
	return bytesToHexSubstring(hash, 4, 12) // 相当于 [8:24]
}

func keccak384HashOptimized(parts ...string) string {
	hasher := keccakPool.Get().(hash.Hash)
	defer func() {
		hasher.Reset()
		keccakPool.Put(hasher)
	}()

	builder := builderPool.Get().(*[]byte)
	defer func() {
		*builder = (*builder)[:0]
		builderPool.Put(builder)
	}()

	totalLen := 0
	for _, part := range parts {
		totalLen += len(part)
	}
	
	if cap(*builder) < totalLen {
		*builder = make([]byte, 0, totalLen)
	}

	for _, part := range parts {
		*builder = append(*builder, part...)
	}

	hasher.Write(*builder)
	hash := hasher.Sum(nil)
	
	return bytesToHexSubstring(hash, 6, 14) // 相当于 [12:28]
}

// 高效的字节到十六进制转换
func bytesToHexSubstring(src []byte, start, end int) string {
	const hexTable = "0123456789abcdef"
	dst := make([]byte, (end-start)*2)
	
	j := 0
	for i := start; i < end && i < len(src); i++ {
		dst[j] = hexTable[src[i]>>4]
		dst[j+1] = hexTable[src[i]&0x0f]
		j += 2
	}
	
	return *(*string)(unsafe.Pointer(&dst)) // 零拷贝转换
}

// 缓存清理
func cleanupCache() {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()
	
	now := time.Now()
	for key, item := range licenseCache {
		if now.Sub(item.timestamp) > cacheTTL {
			delete(licenseCache, key)
		}
	}
	
	// 如果缓存超过最大大小，删除最旧的条目
	if len(licenseCache) > cacheMaxSize {
		oldestKey := ""
		oldestTime := now
		for key, item := range licenseCache {
			if item.timestamp.Before(oldestTime) {
				oldestTime = item.timestamp
				oldestKey = key
			}
		}
		if oldestKey != "" {
			delete(licenseCache, oldestKey)
		}
	}
}

// 保留原始函数用于向后兼容
func md5Hash(msg string) string {
	return md5HashOptimized(msg)
}

func keccak384Hash(msg string) string {
	return keccak384HashOptimized(msg)
}

type FinalShellLicense struct {
	AdvancedBelow396 string `json:"advancedBelow396"`
	ProBelow396      string `json:"proBelow396"`
	AdvancedAbove396 string `json:"advancedAbove396"`
	ProAbove396      string `json:"proAbove396"`
	AdvancedAbove45  string `json:"advancedAbove45"`
	ProAbove45       string `json:"proAbove45"`
	AdvancedAbove46  string `json:"advancedAbove46"`
	ProAbove46       string `json:"proAbove46"`
}

// GenerateLicense generates license codes for different FinalShell versions and editions with caching
func GenerateLicense(machineCode string) FinalShellLicense {
	// 检查缓存
	cacheMutex.RLock()
	if item, exists := licenseCache[machineCode]; exists {
		if time.Since(item.timestamp) < cacheTTL {
			cacheMutex.RUnlock()
			return item.license
		}
	}
	cacheMutex.RUnlock()

	// 批量生成所有许可证 - 使用优化的哈希函数
	license := FinalShellLicense{
		AdvancedBelow396: md5HashOptimized("61305", machineCode, "8552"),
		ProBelow396:      md5HashOptimized("2356", machineCode, "13593"),
		AdvancedAbove396: keccak384HashOptimized(machineCode, "hSf(78cvVlS5E"),
		ProAbove396:      keccak384HashOptimized(machineCode, "FF3Go(*Xvbb5s2"),
		AdvancedAbove45:  keccak384HashOptimized(machineCode, "wcegS3gzA$"),
		ProAbove45:       keccak384HashOptimized(machineCode, "b(xxkHn%z);x"),
		AdvancedAbove46:  keccak384HashOptimized(machineCode, "csSf5*xlkgYSX,y"),
		ProAbove46:       keccak384HashOptimized(machineCode, "Scfg*ZkvJZc,s,Y"),
	}

	// 更新缓存
	cacheMutex.Lock()
	licenseCache[machineCode] = &cacheItem{
		license:   license,
		timestamp: time.Now(),
	}
	
	// 定期清理缓存
	if len(licenseCache)%100 == 0 { // 每100次插入清理一次
		go cleanupCache()
	}
	cacheMutex.Unlock()

	return license
}

// GetCacheStats returns cache statistics
func GetCacheStats() (size int) {
	cacheMutex.RLock()
	defer cacheMutex.RUnlock()
	return len(licenseCache)
}

// ClearCache clears all cached licenses
func ClearCache() {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()
	licenseCache = make(map[string]*cacheItem)
}
