package optimized

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"hash"
	"license/jrebel/constant"
	"license/jrebel/vo"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

// OptimizedLeasesController defines the optimized controller structure with caching and pooling
type OptimizedLeasesController struct {
	// Pre-parsed RSA private key (parsed only once)
	privateKey *rsa.PrivateKey
	
	// Object pools for memory reuse
	stringBuilderPool sync.Pool
	sha1HasherPool    sync.Pool
	
	// Signature cache with TTL
	signatureCache *SignatureCache
	
	// VO object pools
	leasesHandlerVOPool    sync.Pool
	leasesOneHandlerVOPool sync.Pool
	validateHandlerVOPool  sync.Pool
	
	// Pre-allocated byte slice for signature computation
	signatureBuffer []byte
	bufferMutex     sync.Mutex
}

// SignatureCache implements a thread-safe LRU cache with TTL for signatures
type SignatureCache struct {
	cache   map[string]*SignatureCacheEntry
	mutex   sync.RWMutex
	maxSize int
	ttl     time.Duration
}

type SignatureCacheEntry struct {
	signature string
	createdAt time.Time
}

// NewOptimizedLeasesController creates a new optimized controller instance
func NewOptimizedLeasesController() (*OptimizedLeasesController, error) {
	// Pre-parse the RSA private key once
	block, _ := pem.Decode([]byte(constant.LeasesPrivateKey))
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block containing the private key")
	}
	
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse RSA private key: %v", err)
	}
	
	controller := &OptimizedLeasesController{
		privateKey: privateKey,
		signatureCache: &SignatureCache{
			cache:   make(map[string]*SignatureCacheEntry),
			maxSize: 1000, // Configurable cache size
			ttl:     5 * time.Minute, // 5 minute TTL for signatures
		},
		signatureBuffer: make([]byte, 0, 256), // Pre-allocated buffer
	}
	
	// Initialize string builder pool
	controller.stringBuilderPool = sync.Pool{
		New: func() interface{} {
			return &strings.Builder{}
		},
	}
	
	// Initialize SHA1 hasher pool
	controller.sha1HasherPool = sync.Pool{
		New: func() interface{} {
			return sha1.New()
		},
	}
	
	// Initialize VO object pools
	controller.leasesHandlerVOPool = sync.Pool{
		New: func() interface{} {
			return &vo.LeasesHandlerVO{
				ZeroIds: make([]string, 0), // Pre-allocate empty slice
			}
		},
	}
	
	controller.leasesOneHandlerVOPool = sync.Pool{
		New: func() interface{} {
			return &vo.LeasesOneHandlerVO{}
		},
	}
	
	controller.validateHandlerVOPool = sync.Pool{
		New: func() interface{} {
			return &vo.ValidateHandlerVO{}
		},
	}
	
	// Start background cleanup goroutine
	go controller.startCacheCleanup()
	
	return controller, nil
}

// startCacheCleanup runs a background goroutine to cleanup expired cache entries
func (controller *OptimizedLeasesController) startCacheCleanup() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	
	for range ticker.C {
		controller.signatureCache.cleanup()
	}
}

// optimizedSign creates a digital signature using optimized RSA-SHA1 algorithm with caching
func (controller *OptimizedLeasesController) optimizedSign(clientRandomness, guid string, offline bool, validFrom, validUntil int64) string {
	// Create cache key
	cacheKey := fmt.Sprintf("%s:%s:%s:%t:%d:%d", 
		clientRandomness, constant.ServerRandomness, guid, offline, validFrom, validUntil)
	
	// Try to get from cache first
	if cachedSignature := controller.signatureCache.get(cacheKey); cachedSignature != "" {
		return cachedSignature
	}
	
	// Build signature string using pooled string builder
	builder := controller.stringBuilderPool.Get().(*strings.Builder)
	builder.Reset() // Clear the builder for reuse
	defer controller.stringBuilderPool.Put(builder)
	
	// Efficient string building
	builder.WriteString(clientRandomness)
	builder.WriteByte(';')
	builder.WriteString(constant.ServerRandomness)
	builder.WriteByte(';')
	builder.WriteString(guid)
	builder.WriteByte(';')
	builder.WriteString(strconv.FormatBool(offline))
	
	if offline {
		builder.WriteByte(';')
		builder.WriteString(strconv.FormatInt(validFrom, 10))
		builder.WriteByte(';')
		builder.WriteString(strconv.FormatInt(validUntil, 10))
	}
	
	signatureBase := builder.String()
	log.Printf("signature: %s", signatureBase)
	
	// Use pooled SHA1 hasher
	hasher := controller.sha1HasherPool.Get().(hash.Hash)
	hasher.Reset() // Reset for reuse
	defer controller.sha1HasherPool.Put(hasher)
	
	hasher.Write([]byte(signatureBase))
	hashed := hasher.Sum(nil)
	
	// Use pre-parsed private key (no parsing overhead)
	signature, err := rsa.SignPKCS1v15(rand.Reader, controller.privateKey, crypto.SHA1, hashed)
	if err != nil {
		log.Printf("Failed to sign data: %v", err)
		return ""
	}
	
	signatureStr := base64.StdEncoding.EncodeToString(signature)
	
	// Cache the signature
	controller.signatureCache.set(cacheKey, signatureStr)
	
	return signatureStr
}

// Cache methods for SignatureCache
func (sc *SignatureCache) get(key string) string {
	sc.mutex.RLock()
	defer sc.mutex.RUnlock()
	
	if entry, ok := sc.cache[key]; ok {
		if time.Since(entry.createdAt) < sc.ttl {
			return entry.signature
		}
		// Entry expired, will be cleaned up later
	}
	return ""
}

func (sc *SignatureCache) set(key, signature string) {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()
	
	// If cache is full, remove oldest entry
	if len(sc.cache) >= sc.maxSize {
		sc.evictOldest()
	}
	
	sc.cache[key] = &SignatureCacheEntry{
		signature: signature,
		createdAt: time.Now(),
	}
}

func (sc *SignatureCache) evictOldest() {
	var oldestKey string
	var oldestTime time.Time = time.Now()
	
	for key, entry := range sc.cache {
		if entry.createdAt.Before(oldestTime) {
			oldestTime = entry.createdAt
			oldestKey = key
		}
	}
	
	if oldestKey != "" {
		delete(sc.cache, oldestKey)
	}
}

func (sc *SignatureCache) cleanup() {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()
	
	now := time.Now()
	for key, entry := range sc.cache {
		if now.Sub(entry.createdAt) > sc.ttl {
			delete(sc.cache, key)
		}
	}
}

// Optimized LeasesHandler using object pooling and pre-filled structures
func (controller *OptimizedLeasesController) OptimizedLeasesHandler(c *gin.Context) {
	// Parse form data once
	clientRandomness := c.PostForm("randomness")
	username := c.PostForm("username")
	guid := c.PostForm("guid")
	offline, _ := strconv.ParseBool(c.PostForm("offline"))
	clientTime, _ := strconv.ParseInt(c.PostForm("clientTime"), 10, 64)
	
	var validFrom, validUntil int64
	if offline {
		// Calculate time after 180 days, note that the unit is milliseconds
		expiration := clientTime + 180*24*60*60*1000
		validFrom = clientTime
		validUntil = expiration
	}
	
	// Generate signature using optimized method
	signature := controller.optimizedSign(clientRandomness, guid, offline, validFrom, validUntil)
	
	// Get VO object from pool and reset/populate it
	leasesHandlerVO := controller.leasesHandlerVOPool.Get().(*vo.LeasesHandlerVO)
	defer controller.leasesHandlerVOPool.Put(leasesHandlerVO)
	
	// Reset and populate the pooled object
	*leasesHandlerVO = vo.LeasesHandlerVO{
		ServerVersion:         constant.ServerVersion,
		ServerProtocolVersion: constant.ServerProtocolVersion,
		ServerGuid:            constant.ServerGuid,
		GroupType:             constant.GroupType,
		ID:                    1,
		LicenseType:           1,
		EvaluationLicense:     false,
		Signature:             signature,
		ServerRandomness:      constant.ServerRandomness,
		SeatPoolType:          constant.SeatPoolType,
		StatusCode:            constant.StatusCode,
		Offline:               offline,
		ValidFrom:             validFrom,
		ValidUntil:            validUntil,
		Company:               username,
		OrderId:               uuid.NewString(),
		ZeroIds:               leasesHandlerVO.ZeroIds[:0], // Reuse slice, reset length to 0
		LicenseValidFrom:      1490544001000,
		LicenseValidUntil:     4102415999000,
	}
	
	c.JSON(http.StatusOK, leasesHandlerVO)
}

// Optimized Leases1Handler
func (controller *OptimizedLeasesController) OptimizedLeases1Handler(c *gin.Context) {
	username := c.DefaultQuery("username", "")
	
	clientRandomness := c.PostForm("randomness")
	guid := c.PostForm("guid")
	offline, _ := strconv.ParseBool(c.PostForm("offline"))
	clientTime, _ := strconv.ParseInt(c.PostForm("clientTime"), 10, 64)
	
	var validFrom, validUntil int64
	if offline {
		// Calculate time after 180 days, note that the unit is milliseconds
		expiration := clientTime + 180*24*60*60*1000
		validFrom = clientTime
		validUntil = expiration
	}
	
	signature := controller.optimizedSign(clientRandomness, guid, offline, validFrom, validUntil)
	
	// Get VO object from pool
	leasesOneHandlerVO := controller.leasesOneHandlerVOPool.Get().(*vo.LeasesOneHandlerVO)
	defer controller.leasesOneHandlerVOPool.Put(leasesOneHandlerVO)
	
	// Reset and populate the pooled object
	*leasesOneHandlerVO = vo.LeasesOneHandlerVO{
		ServerVersion:         constant.ServerVersion,
		ServerProtocolVersion: constant.ServerProtocolVersion,
		ServerGuid:            constant.ServerGuid,
		Signature:             signature,
		ServerRandomness:      constant.ServerRandomness,
		Features:              "{}",
		GroupType:             constant.GroupType,
		StatusCode:            constant.StatusCode,
		Company:               username,
		Msg:                   "",
		StatusMessage:         "",
	}
	
	c.JSON(http.StatusOK, leasesOneHandlerVO)
}

// Optimized ValidateHandler
func (controller *OptimizedLeasesController) OptimizedValidateHandler(c *gin.Context) {
	clientRandomness := c.PostForm("randomness")
	guid := c.PostForm("guid")
	offline, _ := strconv.ParseBool(c.PostForm("offline"))
	clientTime, _ := strconv.ParseInt(c.PostForm("clientTime"), 10, 64)
	
	var validFrom, validUntil int64
	if offline {
		// Calculate time after 180 days, note that the unit is milliseconds
		expiration := clientTime + 180*24*time.Hour.Milliseconds()
		validFrom = clientTime
		validUntil = expiration
	}
	
	signature := controller.optimizedSign(clientRandomness, guid, offline, validFrom, validUntil)
	
	// Get VO object from pool
	validateHandlerVO := controller.validateHandlerVOPool.Get().(*vo.ValidateHandlerVO)
	defer controller.validateHandlerVOPool.Put(validateHandlerVO)
	
	// Reset and populate the pooled object
	*validateHandlerVO = vo.ValidateHandlerVO{
		ServerVersion:         constant.ServerVersion,
		ServerProtocolVersion: constant.ServerProtocolVersion,
		ServerGuid:            constant.ServerGuid,
		Signature:             signature,
		ServerRandomness:      constant.ServerRandomness,
		Features:              "{}",
		GroupType:             constant.GroupType,
		StatusCode:            constant.StatusCode,
		Company:               constant.Company,
		CanGetLease:           true,
		LicenseType:           "1",
		EvaluationLicense:     false,
		SeatPoolType:          constant.SeatPoolType,
	}
	
	c.JSON(http.StatusOK, validateHandlerVO)
}

// GetCacheStats returns cache statistics for monitoring
func (controller *OptimizedLeasesController) GetCacheStats() map[string]interface{} {
	controller.signatureCache.mutex.RLock()
	defer controller.signatureCache.mutex.RUnlock()
	
	return map[string]interface{}{
		"cache_size": len(controller.signatureCache.cache),
		"max_size":   controller.signatureCache.maxSize,
		"ttl_minutes": controller.signatureCache.ttl.Minutes(),
	}
}

// ClearCache clears the signature cache
func (controller *OptimizedLeasesController) ClearCache() {
	controller.signatureCache.mutex.Lock()
	defer controller.signatureCache.mutex.Unlock()
	
	controller.signatureCache.cache = make(map[string]*SignatureCacheEntry)
}