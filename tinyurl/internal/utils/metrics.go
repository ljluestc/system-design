package utils

import (
    "github.com/sirupsen/logrus"
)

// Metrics tracks system metrics
type Metrics struct {
    log *logrus.Logger
}

// NewMetrics creates a Metrics instance
func NewMetrics(log *logrus.Logger) *Metrics {
    return &Metrics{log: log}
}

// RecordRequest logs a request
func (m *Metrics) RecordRequest() {
    m.log.Info("Recorded request")
}