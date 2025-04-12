package utils

import "github.com/sirupsen/logrus"

func NewLogger() *logrus.Logger {
    log := logrus.New()
    log.SetFormatter(&logrus.TextFormatter{})
    log.SetLevel(logrus.InfoLevel)
    return log
}