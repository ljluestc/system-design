package cache

import (
    "github.com/bradfitz/gomemcache/memcache"
    "github.com/calelin/messenger/config"
    "github.com/sirupsen/logrus"
)

// MemcachedClient wraps Memcached operations
type MemcachedClient struct {
    client *memcache.Client
    log    *logrus.Logger
}

// NewMemcachedClient creates a new MemcachedClient
func NewMemcachedClient(cfg *config.Config, log *logrus.Logger) *MemcachedClient {
    client := memcache.New(cfg.MemcachedHost)
    return &MemcachedClient{client: client, log: log}
}

// Get retrieves a value from cache
func (c *MemcachedClient) Get(key string) (string, error) {
    item, err := c.client.Get(key)
    if err != nil {
        c.log.Debugf("Cache miss for key %s: %v", key, err)
        return "", err
    }
    c.log.Debugf("Cache hit for key %s", key)
    return string(item.Value), nil
}

// Set stores a value in cache
func (c *MemcachedClient) Set(key, value string) error {
    item := &memcache.Item{Key: key, Value: []byte(value), Expiration: 3600} // 1 hour expiry
    if err := c.client.Set(item); err != nil {
        c.log.Errorf("Failed to set cache for key %s: %v", key, err)
        return err
    }
    c.log.Debugf("Cached key %s", key)
    return nil
}

// Delete removes a key from cache
func (c *MemcachedClient) Delete(key string) error {
    if err := c.client.Delete(key); err != nil && err != memcache.ErrCacheMiss {
        c.log.Errorf("Failed to delete cache for key %s: %v", key, err)
        return err
    }
    c.log.Debugf("Deleted cache for key %s", key)
    return nil
}