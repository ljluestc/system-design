package utils

import "time"

// GetCurrentTimestamp returns the current time in milliseconds
func GetCurrentTimestamp() int64 {
    return time.Now().UnixNano() / 1e6
}