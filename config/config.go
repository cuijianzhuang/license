package config

import (
	"fmt"
	"license/logger"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	HttpHost       string
	HttpPort       int
	DataDir        string
	DatabaseDriver string
	DatabaseDsn    string
	StartTime      time.Time
}

var globalConfig *Config

// InitConfig initializes global configuration
func InitConfig() {
	if globalConfig != nil {
		logger.Info("Config is already initialized")
		return
	}
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		logger.Info("No .env file found")
	}

	dataDir := getEnvStr("DATA_DIR", "/data")
	databaseDriver := getEnvStr("DATABASE_DRIVER", "sqlite")
	var databaseDsn string
	if databaseDriver == "sqlite" {
		databaseDsn = fmt.Sprintf("%s/%s", dataDir, getEnvStr("DATABASE_DSN", "license.db"))
	} else {
		databaseDsn = getEnvStr("DATABASE_DSN", "")
	}

	globalConfig = &Config{
		HttpHost:       getEnvStr("HTTP_HOST", "0.0.0.0"),
		HttpPort:       getEnvInt("HTTP_PORT", 5000),
		DataDir:        dataDir,
		DatabaseDriver: databaseDriver,
		DatabaseDsn:    databaseDsn,
		StartTime:      time.Now(),
	}
}

// getEnvStr reads environment variable or returns default value
func getEnvStr(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// getEnvInt reads environment variable or returns default value
func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		intValue, err := strconv.Atoi(value)
		if err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvBool reads environment variable or returns default value
func getEnvBool(key string, defaultValue bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		boolValue, err := strconv.ParseBool(value)
		if err == nil {
			return boolValue
		}
	}
	return defaultValue
}

// GetConfig provides access to global configuration
func GetConfig() *Config {
	if globalConfig == nil {
		logger.Info("Config is not initialized")
		// Not initialized, perform initialization
		InitConfig()
	}
	return globalConfig
}
