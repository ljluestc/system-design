package utils

import (
    "github.com/sirupsen/logrus"
)

// NewLogger creates a logger
func NewLogger() *logrus.Logger {
    log := logrus.New()
    log.SetFormatter(&logrus.TextFormatter{})
    log.SetLevel(logrus.InfoLevel)
    return log
}