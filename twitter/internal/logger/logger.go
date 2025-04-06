package logger

import (
    "log"
    "os"
)

// Logger wraps the standard log package with custom formatting
type Logger struct {
    *log.Logger
}

// NewLogger initializes a new logger instance
func NewLogger() *Logger {
    return &Logger{log.New(os.Stdout, "TWITTER: ", log.LstdFlags)}
}

// Info logs informational messages
func (l *Logger) Info(msg string, args ...interface{}) {
    l.Printf("[INFO] "+msg, args...)
}

// Error logs error messages
func (l *Logger) Error(msg string, args ...interface{}) {
    l.Printf("[ERROR] "+msg, args...)
}