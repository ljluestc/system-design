package utils

import (
    "crypto/sha256"
    "fmt"
)

// HashKey generates a SHA-256 hash
func HashKey(key string) string {
    hash := sha256.Sum256([]byte(key))
    return fmt.Sprintf("%x", hash)
}

// Expand with more utilities
// ...