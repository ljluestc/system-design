package tests

import (
    "github.com/calelin/messenger/internal/cache"
    "github.com/calelin/messenger/config"
    "github.com/sirupsen/logrus"
    "testing"
)

func TestCacheSetGet(t *testing.T) {
    log := logrus.New()
    cfg := config.NewConfig()
    client := cache.NewMemcachedClient(cfg, log)
    key := "test-key"
    value := "test-value"

    if err := client.Set(key, value); err != nil {
        t.Errorf("Failed to set cache: %v", err)
    }

    result, err := client.Get(key)
    if err != nil {
        t.Errorf("Failed to get cache: %v", err)
    }
    if result != value {
        t.Errorf("Expected %s, got %s", value, result)
    }
}