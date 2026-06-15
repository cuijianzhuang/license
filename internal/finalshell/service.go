package finalshell

import (
	"encoding/hex"
	"hash"
	"sync"
	"time"

	"crypto/md5"

	"github.com/ebfe/keccak"
)

const (
	cacheMaxSize = 1000
	cacheTTL     = 30 * time.Minute
)

type cacheItem struct {
	license   FinalShellLicense
	timestamp time.Time
}

var (
	licenseCache = make(map[string]*cacheItem)
	cacheMutex   sync.RWMutex
)

func md5Hash(parts ...string) string {
	return digestHex(md5.New(), parts, 4, 12)
}

func keccak384Hash(parts ...string) string {
	return digestHex(keccak.New384(), parts, 6, 14)
}

// digestHex writes parts into h, hex-encodes the digest, then returns the
// substring [start*2 : end*2]. The two-step slice keeps the original
// "[8:24]" / "[12:28]" hex window selection intact.
func digestHex(h hash.Hash, parts []string, start, end int) string {
	for _, p := range parts {
		h.Write([]byte(p))
	}
	full := hex.EncodeToString(h.Sum(nil))
	return full[start*2 : end*2]
}

func cleanupCache() {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	now := time.Now()
	for key, item := range licenseCache {
		if now.Sub(item.timestamp) > cacheTTL {
			delete(licenseCache, key)
		}
	}

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

// GenerateLicense returns the eight license codes for a machine fingerprint,
// memoised per machineCode for cacheTTL.
func GenerateLicense(machineCode string) FinalShellLicense {
	cacheMutex.RLock()
	if item, exists := licenseCache[machineCode]; exists && time.Since(item.timestamp) < cacheTTL {
		cacheMutex.RUnlock()
		return item.license
	}
	cacheMutex.RUnlock()

	license := FinalShellLicense{
		AdvancedBelow396: md5Hash("61305", machineCode, "8552"),
		ProBelow396:      md5Hash("2356", machineCode, "13593"),
		AdvancedAbove396: keccak384Hash(machineCode, "hSf(78cvVlS5E"),
		ProAbove396:      keccak384Hash(machineCode, "FF3Go(*Xvbb5s2"),
		AdvancedAbove45:  keccak384Hash(machineCode, "wcegS3gzA$"),
		ProAbove45:       keccak384Hash(machineCode, "b(xxkHn%z);x"),
		AdvancedAbove46:  keccak384Hash(machineCode, "csSf5*xlkgYSX,y"),
		ProAbove46:       keccak384Hash(machineCode, "Scfg*ZkvJZc,s,Y"),
	}

	cacheMutex.Lock()
	licenseCache[machineCode] = &cacheItem{
		license:   license,
		timestamp: time.Now(),
	}
	needsCleanup := len(licenseCache)%100 == 0
	cacheMutex.Unlock()

	if needsCleanup {
		go cleanupCache()
	}

	return license
}

// GetCacheStats returns the current cached entry count.
func GetCacheStats() int {
	cacheMutex.RLock()
	defer cacheMutex.RUnlock()
	return len(licenseCache)
}

// ClearCache drops all cached licenses.
func ClearCache() {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()
	licenseCache = make(map[string]*cacheItem)
}
