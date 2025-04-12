package cache

import (
    "sync"
    "time"
)

type CacheNode struct {
    mu       sync.RWMutex   // Mutex for thread safety
    data     map[string]*CacheEntry
    eviction *LRUEviction   // Eviction policy
}

type CacheEntry struct {
    Value      string    // Cached value
    Expiration time.Time // Expiration time
}

func NewCacheNode(capacity int) *CacheNode {
    return &CacheNode{
        data:     make(map[string]*CacheEntry),
        eviction: NewLRUEviction(capacity),
    }
}

func (n *CacheNode) Set(key, value string, ttl time.Duration) {
    n.mu.Lock()
    defer n.mu.Unlock()
    // Evict if capacity is reached and key doesn't exist
    if _, exists := n.data[key]; !exists && len(n.data) >= n.eviction.capacity {
        oldest := n.eviction.list.Back()
        if oldest != nil {
            oldKey := oldest.Value.(string)
            delete(n.data, oldKey)
            n.eviction.Remove(oldKey)
        }
    }
    expiration := time.Now().Add(ttl)
    n.data[key] = &CacheEntry{Value: value, Expiration: expiration}
    n.eviction.Touch(key)
}

func (n *CacheNode) Get(key string) (string, bool) {
    n.mu.RLock()
    defer n.mu.RUnlock()
    entry, exists := n.data[key]
    if !exists || entry.Expiration.Before(time.Now()) {
        return "", false
    }
    n.eviction.Touch(key)
    return entry.Value, true
}