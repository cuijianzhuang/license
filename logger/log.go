package logger

import (
	"sync"

	"github.com/sirupsen/logrus"
)

var (
	instance *LogWrapper
	once     sync.Once
)

// customFormatter custom log format
type customFormatter struct{}

// Format constructs the log output format, ensuring there are no extra prefixes or field names
func (f *customFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := entry.Time.Format("2006/01/02 - 15:04:05")
	return []byte("[" + entry.Data["level"].(string) + "] " + timestamp + " " + entry.Message + "\n"), nil
}

// LogWrapper wraps the logrus logger
type LogWrapper struct {
	*logrus.Logger
}

// initLogger creates and initializes the logger instance
func initLogger() {
	baseLogger := logrus.New()
	baseLogger.SetFormatter(new(customFormatter)) // Use custom format
	instance = &LogWrapper{baseLogger}
}

// GetInstance returns the singleton logger instance
func GetInstance() *LogWrapper {
	once.Do(initLogger)
	return instance
}

// Info logs an INFO level message
func Info(message string) {
	GetInstance().Logger.WithField("level", "INFO").Info(message)
}

// Error logs an ERROR level message with an error object
func Error(message string, err error) {
	if err != nil {
		GetInstance().Logger.WithField("level", "ERROR").Errorf("%s | Error: %v", message, err)
	} else {
		GetInstance().Logger.WithField("level", "ERROR").Error(message)
	}
}

// Sys logs a SYS level message
func Sys(message string) {
	GetInstance().Logger.WithField("level", "SYS").Info(message)
}
