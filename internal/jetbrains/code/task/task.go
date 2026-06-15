package task

import (
	"license/internal/jetbrains/code/service"
	"license/internal/logger"
)

// FetchProductLatest refreshes products and plugins from the upstream catalog.
func FetchProductLatest() {
	if err := service.FetchLatestProducts(); err != nil {
		logger.Error("Failed to fetch latest product:", err)
	}
	if err := service.FetchLatestPlugins(); err != nil {
		logger.Error("Failed to fetch latest plugin:", err)
	}
}
