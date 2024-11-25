package logger

import (
	"PaymentAPI/constants"
	"os"

	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func init() {
	// Initialize the logger
	Logger = logrus.New()

	// Open or create the log.json file
	file, err := os.OpenFile(constants.LogJsonPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logrus.Fatalf("Failed to open log file: %v", err)
	}

	// Configure the logger
	Logger.SetOutput(file)
	Logger.SetFormatter(&logrus.JSONFormatter{})
	Logger.SetLevel(logrus.InfoLevel)
}

// LogInfo logs an informational message
func LogInfo(message string, fields logrus.Fields) {
	Logger.WithFields(fields).Info(message)
}

// LogError logs an error message
func LogError(message string, fields logrus.Fields) {
	Logger.WithFields(fields).Error(message)
}

// LogWarning logs a warning message
func LogWarning(message string, fields logrus.Fields) {
	Logger.WithFields(fields).Warn(message)
}
