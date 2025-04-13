package dns

import (
    "net"
    "time"
    "crawler/internal/utils"
)

type Resolver struct {
    cache *Cache
}

func New(cacheTTL time.Duration) *Resolver {
    return &Resolver{cache: NewCache(cacheTTL)}
}

func (r *Resolver) Resolve(hostname string) (string, error) {
    // Check cache
    if ip, ok := r.cache.Get(hostname); ok {
        return ip, nil
    }

    // Perform DNS lookup
    addrs, err := net.LookupHost(hostname)
    if err != nil {
        return "", err
    }
    if len(addrs) == 0 {
        return "", utils.Errorf("no IPs for %s", hostname)
    }

    // Cache result
    r.cache.Set(hostname, addrs[0])
    return addrs[0], nil
}