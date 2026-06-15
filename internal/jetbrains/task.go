package jetbrains

import (
	"license/internal/logger"
)

// FetchProductLatest refreshes products and plugins from the upstream catalog.
func FetchProductLatest() {
	if err := FetchLatestProducts(); err != nil {
		logger.Error("Failed to fetch latest product:", err)
	}
	if err := FetchLatestPlugins(); err != nil {
		logger.Error("Failed to fetch latest plugin:", err)
	}
}
