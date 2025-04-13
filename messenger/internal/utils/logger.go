package utils

import (
    "log"
    "os"
)

type Logger struct {
    *log.Logger
}

func NewLogger() *Logger {
    return &Logger{log.New(os.Stdout, "messenger: ", log.LstdFlags|log.Lshortfile)}
}

func (l *Logger) Info(msg string) {
    l.Printf("INFO: %s", msg)
}

func (l *Logger) Error(msg string, err error) {
    l.Printf("ERROR: %s: %v", msg, err)
}