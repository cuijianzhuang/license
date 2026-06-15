package logger

import (
	"github.com/sirupsen/logrus"
)

// customFormatter renders logrus entries as "[LEVEL] yyyy/mm/dd - HH:MM:SS msg".
type customFormatter struct{}

// Format constructs the log output format, ensuring there are no extra prefixes or field names
func (f *customFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := entry.Time.Format("2006/01/02 - 15:04:05")
	return []byte("[" + entry.Data["level"].(string) + "] " + timestamp + " " + entry.Message + "\n"), nil
}

// std is the package-level logrus instance shared by all helpers below. It is
// initialised at package load and never reassigned, so no mutex is needed.
var std = func() *logrus.Logger {
	l := logrus.New()
	l.SetFormatter(new(customFormatter))
	return l
}()

// Info logs an INFO level message
func Info(message string) {
	std.WithField("level", "INFO").Info(message)
}

// Error logs an ERROR level message with an error object
func Error(message string, err error) {
	if err != nil {
		std.WithField("level", "ERROR").Errorf("%s | Error: %v", message, err)
	} else {
		std.WithField("level", "ERROR").Error(message)
	}
}

// Sys logs a SYS level message
func Sys(message string) {
	std.WithField("level", "SYS").Info(message)
}
