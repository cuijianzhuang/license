package v2

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"license/internal/config"
	"license/internal/jetbrains/types"
	"license/internal/logger"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/net/proxy"
	"gorm.io/gorm/clause"
)

// getHTTPClient returns an HTTP client with proxy support from config
// Proxy priority (highest to lowest):
// 1. ALL_PROXY - SOCKS5 proxy (socks5://, socks4://)
// 2. HTTPS_PROXY - HTTP proxy for HTTPS requests
// 3. HTTP_PROXY - HTTP proxy for all requests
// 4. Direct connection
func getHTTPClient() *http.Client {
	cfg := config.GetConfig()
	transport := &http.Transport{}

	var usingProxy string

	// Check for SOCKS proxy via ALL_PROXY
	if cfg.ALLProxy != "" && (strings.HasPrefix(strings.ToLower(cfg.ALLProxy), "socks5") ||
		strings.HasPrefix(strings.ToLower(cfg.ALLProxy), "socks4") ||
		strings.HasPrefix(strings.ToLower(cfg.ALLProxy), "socks://")) {
		parsed, err := url.Parse(cfg.ALLProxy)
		if err == nil {
			dialer, err := proxy.FromURL(parsed, proxy.Direct)
			if err == nil {
				transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
					return dialer.Dial(network, addr)
				}
				usingProxy = fmt.Sprintf("SOCKS(%s)", cfg.ALLProxy)
			}
		}
	}

	// For HTTP proxy, set Proxy function explicitly
	if usingProxy == "" {
		proxyStr := cfg.HTTPSProxy
		if proxyStr == "" {
			proxyStr = cfg.HTTPProxy
		}
		if proxyStr != "" {
			parsed, err := url.Parse(proxyStr)
			if err == nil {
				transport.Proxy = http.ProxyURL(parsed)
				usingProxy = proxyStr
			}
		}
	}

	// Log proxy configuration
	logger.Info(fmt.Sprintf("Proxy config - HTTP_PROXY: %s, HTTPS_PROXY: %s, ALL_PROXY: %s, Using: %s",
		getOrNone(cfg.HTTPProxy), getOrNone(cfg.HTTPSProxy), getOrNone(cfg.ALLProxy), getOrNone(usingProxy)))

	return &http.Client{
		Transport: transport,
	}
}

func getOrNone(s string) string {
	if s == "" {
		return "(none)"
	}
	return s
}

// ProductService handles product-related operations.
//
// GORM is safe for concurrent use, so the service holds no state and no lock —
// an earlier RWMutex caused all reads to block while FetchLatest held the write
// lock for the duration of a multi-minute HTTP scrape.
type ProductService struct{}

// NewProductService creates a new product service
func NewProductService() *ProductService {
	return &ProductService{}
}

func listProducts() ([]types.Product, error) {
	var products []types.Product
	if err := config.DB.Find(&products).Error; err != nil {
		logger.Error("Failed to fetch products:", err)
		return nil, err
	}
	return products, nil
}

func upsertProducts(products []*types.Product) error {
	err := config.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "product_code"}},
		DoUpdates: clause.AssignmentColumns([]string{"product_name", "product_detail"}),
	}).CreateInBatches(products, 100).Error
	if err != nil {
		logger.Error("Failed to upsert products:", err)
	}
	return err
}

// GetAll retrieves all products
func (s *ProductService) GetAll() ([]types.Product, error) {
	products, err := listProducts()
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
	}
	return products, nil
}

// GetByCode retrieves a product by its code
func (s *ProductService) GetByCode(code string) (*types.Product, error) {
	products, err := listProducts()
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
	}

	for _, product := range products {
		if product.ProductCode == code {
			return &product, nil
		}
	}

	return nil, fmt.Errorf("product not found: %s", code)
}

// FetchLatest fetches the latest products from external source
func (s *ProductService) FetchLatest() error {
	client := getHTTPClient()
	req, err := http.NewRequest("GET", "https://data.services.jetbrains.com/products", nil)
	if err != nil {
		logger.Error("Error creating request:", err)
		return err
	}

	req.Header.Set("User-Agent", getUserAgent())
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("Error executing request:", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Error(fmt.Sprintf("Failed to fetch product information with status code: %d", resp.StatusCode), nil)
		return fmt.Errorf("failed to fetch product information with status code: %d", resp.StatusCode)
	}

	var products []map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&products)
	if err != nil {
		logger.Error("Error decoding JSON:", err)
		return err
	}

	productList := make([]*types.Product, 0, len(products))
	for i, product := range products {
		logger.Info(fmt.Sprintf("Total products to process: %d, currently processing #%d", len(products), i+1))

		// Convert product map to JSON string
		productJSON, err := json.Marshal(product)
		if err != nil {
			logger.Error("Error marshaling product to JSON:", err)
			continue
		}

		productEntity := &types.Product{
			ProductDetail: string(productJSON),
			ProductCode:   fmt.Sprint(product["code"]),
			ProductName:   fmt.Sprint(product["name"]),
		}
		productList = append(productList, productEntity)
		// Simulate pause
		time.Sleep(time.Duration(100+rand.Intn(400)) * time.Millisecond)
	}

	if len(productList) > 0 {
		if err := upsertProducts(productList); err != nil {
			return err
		}
	}

	return nil
}

// PluginService handles plugin-related operations.
// See ProductService for why no lock is held.
type PluginService struct{}

// NewPluginService creates a new plugin service
func NewPluginService() *PluginService {
	return &PluginService{}
}

func listPlugins() ([]types.Plugin, error) {
	var plugins []types.Plugin
	if err := config.DB.Find(&plugins).Error; err != nil {
		logger.Error("Failed to fetch plugins:", err)
		return nil, err
	}
	return plugins, nil
}

func upsertPlugins(plugins []*types.Plugin) error {
	err := config.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "plugin_code"}},
		DoUpdates: clause.AssignmentColumns([]string{"plugin_id", "plugin_name", "plugin_api_detail"}),
	}).CreateInBatches(plugins, 100).Error
	if err != nil {
		logger.Error("Failed to upsert plugins:", err)
	}
	return err
}

// GetAll retrieves all plugins
func (s *PluginService) GetAll() ([]types.Plugin, error) {
	plugins, err := listPlugins()
	if err != nil {
		return nil, fmt.Errorf("failed to get plugins: %w", err)
	}
	return plugins, nil
}

const (
	pluginsBaseURL  = "https://plugins.jetbrains.com/api/searchPlugins?excludeTags=theme&max=24&offset=%d&orderBy=downloads&pricingModels=%s"
	pluginDetailURL = "https://plugins.jetbrains.com/api/plugins/"
	maxPerPage      = 24
)

var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.2 Safari/605.1.15",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:133.0) Gecko/20100101 Firefox/133.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:133.0) Gecko/20100101 Firefox/133.0",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Edge/131.0.0.0",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64; rv:132.0) Gecko/20100101 Firefox/132.0",
}

func getUserAgent() string {
	return userAgents[rand.Intn(len(userAgents))]
}

// GetByCode retrieves a plugin by its code
func (s *PluginService) GetByCode(code string) (*types.Plugin, error) {
	plugins, err := listPlugins()
	if err != nil {
		return nil, fmt.Errorf("failed to get plugins: %w", err)
	}

	for _, plugin := range plugins {
		if plugin.PluginCode == code {
			return &plugin, nil
		}
	}

	return nil, fmt.Errorf("plugin not found: %s", code)
}

// fetchPlugins fetches plugins from external source with pagination
func (s *PluginService) fetchPlugins(pricingModel string) ([]*types.Plugin, error) {
	client := getHTTPClient()

	// Phase 1: Fetch all plugin IDs with pagination
	type pluginInfo struct {
		ID   uint64
		Name string
	}
	var allPluginInfos []pluginInfo
	offset := 0

	for {
		url := fmt.Sprintf(pluginsBaseURL, offset, pricingModel)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			logger.Error("Error creating request:", err)
			return nil, err
		}

		req.Header.Set("User-Agent", getUserAgent())
		resp, err := client.Do(req)
		if err != nil {
			logger.Error("Error on request:", err)
			return nil, err
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			logger.Error(fmt.Sprintf("Failed to fetch plugins, status code: %d", resp.StatusCode), nil)
			return nil, fmt.Errorf("failed to fetch plugins, status code: %d", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			logger.Error("Error reading response body:", err)
			return nil, err
		}

		var data struct {
			Total   int64 `json:"total"`
			Plugins []struct {
				ID   uint64 `json:"id"`
				Name string `json:"name"`
			}
		}

		err = json.Unmarshal(body, &data)
		if err != nil {
			logger.Error("Error unmarshaling JSON:", err)
			return nil, err
		}

		if len(data.Plugins) == 0 {
			break
		}

		for _, p := range data.Plugins {
			allPluginInfos = append(allPluginInfos, pluginInfo{
				ID:   p.ID,
				Name: p.Name,
			})
		}

		logger.Info(fmt.Sprintf("Fetched %d plugins (total: %d, offset: %d, pricingModel: %s)", len(data.Plugins), data.Total, offset, pricingModel))

		if int64(offset+maxPerPage) >= data.Total {
			break
		}

		offset += maxPerPage
	}

	// Phase 2: Fetch details concurrently using worker pool
	const maxWorkers = 10
	pluginCh := make(chan pluginInfo, len(allPluginInfos))
	resultCh := make(chan *types.Plugin, len(allPluginInfos))

	// Send all plugin infos to workers
	for _, p := range allPluginInfos {
		pluginCh <- p
	}
	close(pluginCh)

	var wg sync.WaitGroup
	processed := &atomic.Int64{}

	// Start worker pool
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for p := range pluginCh {
				count := processed.Add(1)
				logger.Info(fmt.Sprintf("Total plugins to process: %d, currently processing #%d, Plugin ID: %d", len(allPluginInfos), count, p.ID))

				detailReq, err := http.NewRequest("GET", fmt.Sprintf("%s%d", pluginDetailURL, p.ID), nil)
				if err != nil {
					logger.Error(fmt.Sprintf("Error creating detail request for ID %d", p.ID), err)
					continue
				}
				detailReq.Header.Set("User-Agent", getUserAgent())

				detailResp, err := client.Do(detailReq)
				if err != nil {
					logger.Error(fmt.Sprintf("Error fetching plugin detail for ID %d", p.ID), err)
					continue
				}

				if detailResp.StatusCode != http.StatusOK {
					detailResp.Body.Close()
					logger.Error(fmt.Sprintf("Failed to fetch plugin detail for ID %d, status: %d", p.ID, detailResp.StatusCode), nil)
					continue
				}

				detailBody, err := io.ReadAll(detailResp.Body)
				detailResp.Body.Close()
				if err != nil {
					logger.Error("Error reading plugin detail response:", err)
					continue
				}

				var detail struct {
					Name         string `json:"name"`
					PurchaseInfo struct {
						ProductCode string `json:"productCode"`
					} `json:"purchaseInfo"`
				}

				err = json.Unmarshal(detailBody, &detail)
				if err != nil {
					logger.Error("Error unmarshaling plugin detail JSON:", err)
					continue
				}

				resultCh <- &types.Plugin{
					PluginID:        p.ID,
					PluginName:      detail.Name,
					PluginCode:      detail.PurchaseInfo.ProductCode,
					PluginApiDetail: string(detailBody),
				}

				// Reduced pause for rate limiting
				time.Sleep(time.Duration(50+rand.Intn(100)) * time.Millisecond)
			}
		}(i)
	}

	// Close result channel when all workers finish
	go func() {
		wg.Wait()
		close(resultCh)
	}()

	// Collect results
	allPlugins := make([]*types.Plugin, 0, len(allPluginInfos))
	for plugin := range resultCh {
		allPlugins = append(allPlugins, plugin)
	}

	return allPlugins, nil
}

// FetchLatest fetches the latest plugins from external source
func (s *PluginService) FetchLatest() error {
	// First fetch paid plugins
	paidPlugins, err := s.fetchPlugins("PAID")
	if err != nil {
		logger.Error("Error fetching paid plugins:", err)
		return err
	}

	// Then fetch freemium plugins
	freemiumPlugins, err := s.fetchPlugins("FREEMIUM")
	if err != nil {
		logger.Error("Error fetching freemium plugins:", err)
		return err
	}

	// Merge the two results
	allPlugins := append(paidPlugins, freemiumPlugins...)

	// Upsert plugins (insert new, update existing)
	if len(allPlugins) > 0 {
		if err := upsertPlugins(allPlugins); err != nil {
			logger.Error("Error upserting plugin batch:", err)
			return err
		}
	}

	return nil
}

