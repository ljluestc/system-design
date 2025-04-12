package utils

import (
    "github.com/google/uuid"
)

// GenerateUUID creates a new UUID
func GenerateUUID() string {
    return uuid.New().String()
}

// Expand with UUID utilities
// ...