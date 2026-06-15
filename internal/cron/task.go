package cron

import (
	"license/internal/jetbrains"
	"license/internal/logger"

	"github.com/robfig/cron/v3"
)

func InitCron() {
	c := cron.New(cron.WithSeconds())

	_, err := c.AddFunc("0 0 1 * * ?", jetbrains.FetchProductLatest)
	if err != nil {
		logger.Error("Failed to add cron job:", err)
	}
	c.Start()
	logger.Sys("Cron job started")
}
