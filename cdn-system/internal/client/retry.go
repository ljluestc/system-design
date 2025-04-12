package client

import (
    "errors"
    "log"
    "time"
)

// Retry executes a function with retries
func Retry(fn func() error, attempts int) error {
    for i := 0; i < attempts; i++ {
        if err := fn(); err == nil {
            return nil
        }
        log.Printf("Retry attempt %d failed", i+1)
        time.Sleep(time.Second)
    }
    return errors.New("max retries exceeded")
}

// Expand with backoff logic
// ...