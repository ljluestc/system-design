package edge

import (
    "log"
    "sync"
    "time"
)

// CacheEntry holds a cached item
type CacheEntry struct {
    Value      string
    Expiration time.Time
}

// Cache manages in-memory content caching
type Cache struct {
    store map[string]CacheEntry
    mu    sync.Mutex
    ttl   time.Duration
}

// NewCache initializes a cache with specified TTL
func NewCache(ttl time.Duration) *Cache {
    return &Cache{
        store: make(map[string]CacheEntry),
        ttl:   ttl,
    }
}

// Get retrieves a cached value if not expired
func (c *Cache) Get(key string) (string, bool) {
    c.mu.Lock()
    defer c.mu.Unlock()

    entry, ok := c.store[key]
    if !ok || entry.Expiration.Before(time.Now()) {
        if ok {
            delete(c.store, key)
            log.Printf("Removed expired key %s", key)
        }
        return "", false
    }
    return entry.Value, true
}

// Set caches a value with TTL
func (c *Cache) Set(key, value string) {
    c.mu.Lock()
    defer c.mu.Unlock()

    c.store[key] = CacheEntry{
        Value:      value,
        Expiration: time.Now().Add(c.ttl),
    }
}

// Invalidate removes a specific key
func (c *Cache) Invalidate(key string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    delete(c.store, key)
}

// Clear resets the cache
func (c *Cache) Clear() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.store = make(map[string]CacheEntry)
}

// Add more cache management functions
// ...