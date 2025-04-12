package edge

import (
    "errors"
    "fmt"
    "log"
    "sync"
    "time"

    "cdn-system/internal/config"
    "cdn-system/internal/origin"
)

// EdgeServer represents an edge server instance
type EdgeServer struct {
    ID         string
    Cache      *Cache
    LastHealth time.Time
    mu         sync.Mutex
}

// Global edge server list
var (
    servers []*EdgeServer
    mu      sync.Mutex
)

// InitServers initializes edge servers based on config
func InitServers(cfg *config.Config) error {
    mu.Lock()
    defer mu.Unlock()
    if len(servers) > 0 {
        return errors.New("edge servers already initialized")
    }
    for i := 0; i < cfg.EdgeServers; i++ {
        servers = append(servers, &EdgeServer{
            ID:         fmt.Sprintf("edge-%d", i),
            Cache:      NewCache(time.Duration(cfg.CacheTTL) * time.Second),
            LastHealth: time.Now(),
        })
    }
    log.Printf("Initialized %d edge servers", cfg.EdgeServers)
    return nil
}

// GetContent fetches content from cache or origin
func GetContent(key string) (string, error) {
    if len(servers) == 0 {
        return "", errors.New("no edge servers available")
    }
    server := servers[0] // Simplified selection
    server.mu.Lock()
    defer server.mu.Unlock()

    if content, ok := server.Cache.Get(key); ok {
        log.Printf("Cache hit for key %s on %s", key, server.ID)
        return content, nil
    }

    content, err := origin.GetContent(key)
    if err != nil {
        log.Printf("Origin fetch failed for key %s: %v", key, err)
        return "", err
    }

    server.Cache.Set(key, content)
    log.Printf("Cached content for key %s on %s", key, server.ID)
    return content, nil
}

// CheckHealth verifies all edge servers are healthy
func CheckHealth() bool {
    mu.Lock()
    defer mu.Unlock()
    for _, server := range servers {
        if !server.IsHealthy() {
            return false
        }
    }
    return true
}

// InvalidateCache clears cache for a specific key
func InvalidateCache(key string) error {
    mu.Lock()
    defer mu.Unlock()
    for _, server := range servers {
        server.mu.Lock()
        server.Cache.Invalidate(key)
        server.mu.Unlock()
        log.Printf("Invalidated cache for key %s on %s", key, server.ID)
    }
    return nil
}

// Expand with additional logic
// ...