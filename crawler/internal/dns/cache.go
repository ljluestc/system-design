package dns

import (
    "sync"
    "time"
)

type Cache struct {
    entries map[string]entry
    ttl     time.Duration
    mu      sync.RWMutex
}

type entry struct {
    ip        string
    expiresAt time.Time
}

func NewCache(ttl time.Duration) *Cache {
    return &Cache{
        entries: make(map[string]entry),
        ttl:     ttl,
    }
}

func (c *Cache) Get(hostname string) (string, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    e, ok := c.entries[hostname]
    if !ok || time.Now().After(e.expiresAt) {
        return "", false
    }
    return e.ip, true
}

func (c *Cache) Set(hostname, ip string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.entries[hostname] = entry{
        ip:        ip,
        expiresAt: time.Now().Add(c.ttl),
    }
}