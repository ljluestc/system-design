package origin

import (
    "errors"
    "fmt"
    "log"
)

// UpdateContentWithVersion updates content with a version
func UpdateContentWithVersion(key, value string, version int) error {
    if server == nil {
        return errors.New("origin server not initialized")
    }
    server.mu.Lock()
    defer server.mu.Unlock()
    // Simplified versioning logic
    server.Content[fmt.Sprintf("%s-v%d", key, version)] = value
    log.Printf("Updated version %d for key %s", version, key)
    return nil
}

// GetContentWithVersion retrieves versioned content
func GetContentWithVersion(key string, version int) (string, error) {
    if server == nil {
        return "", errors.New("origin server not initialized")
    }
    server.mu.Lock()
    defer server.mu.Unlock()
    content, ok := server.Content[fmt.Sprintf("%s-v%d", key, version)]
    if !ok {
        return "", errors.New("version not found")
    }
    return content, nil
}

// Expand with versioning logic
// ...