package config

import (
	"fmt"
	"license/internal/logger"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func SetupDatabase() {
	var err error

	driver := GetConfig().DatabaseDriver
	dsn := GetConfig().DatabaseDsn
	if driver == "sqlite" {
		DB, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
		if err != nil {
			logger.Error("Failed to connect to database:", err)
		}
	} else if driver == "mysql" {
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			logger.Error("Failed to connect to database:", err)
		}
	} else if driver == "postgres" {
		DB, err = gorm.Open(postgres.New(postgres.Config{
			DSN:                  dsn,
			PreferSimpleProtocol: true, // disables implicit prepared statement usage
			DriverName:           "pgx",
		}), &gorm.Config{})
		if err != nil {
			logger.Error("Failed to connect to database:", err)
		}
	} else {
		logger.Error("Unsupported database driver", fmt.Errorf("unsupported driver: %s", driver))
		return
	}

	logger.Sys("Database Connected Successfully")
}
