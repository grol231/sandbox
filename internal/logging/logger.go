package logging

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

// Init initializes the logger with specified level and format
func Init(level, format string) *logrus.Logger {
	log = logrus.New()
	
	// Set output to stdout for container compatibility
	log.SetOutput(os.Stdout)

	// Set log level
	switch level {
	case "debug":
		log.SetLevel(logrus.DebugLevel)
	case "info":
		log.SetLevel(logrus.InfoLevel)
	case "warn":
		log.SetLevel(logrus.WarnLevel)
	case "error":
		log.SetLevel(logrus.ErrorLevel)
	default:
		log.SetLevel(logrus.InfoLevel)
	}

	// Set format for Loki parsing
	if format == "json" {
		log.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
		})
	} else {
		log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
		})
	}

	return log
}

// GetLogger returns the global logger instance
func GetLogger() *logrus.Logger {
	if log == nil {
		return Init("info", "json")
	}
	return log
}

// Info logs an info message
func Info(msg string, fields ...logrus.Fields) {
	entry := GetLogger().WithTime(time.Now())
	if len(fields) > 0 {
		entry = entry.WithFields(fields[0])
	}
	entry.Info(msg)
}

// Error logs an error message
func Error(msg string, err error, fields ...logrus.Fields) {
	entry := GetLogger().WithTime(time.Now())
	if err != nil {
		entry = entry.WithError(err)
	}
	if len(fields) > 0 {
		entry = entry.WithFields(fields[0])
	}
	entry.Error(msg)
}

// Debug logs a debug message
func Debug(msg string, fields ...logrus.Fields) {
	entry := GetLogger().WithTime(time.Now())
	if len(fields) > 0 {
		entry = entry.WithFields(fields[0])
	}
	entry.Debug(msg)
}

// Warn logs a warning message
func Warn(msg string, fields ...logrus.Fields) {
	entry := GetLogger().WithTime(time.Now())
	if len(fields) > 0 {
		entry = entry.WithFields(fields[0])
	}
	entry.Warn(msg)
}