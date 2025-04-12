package client

import (
    "errors"
    "time"
)

func Retry(fn func() error, attempts int) error {
    for i := 0; i < attempts; i++ {
        if err := fn(); err == nil {
            return nil
        }
        time.Sleep(time.Second)
    }
    return errors.New("max retries exceeded")
}