package main

import (
    "context"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
)

func TestCrawler(t *testing.T) {
    crawler, err := NewCrawler()
    assert.NoError(t, err)

    ctx := context.Background()

    t.Run("Enqueue and Crawl", func(t *testing.T) {
        err := crawler.EnqueueURL(ctx, "https://example.com")
        assert.NoError(t, err)

        // Simulate crawl (run for a short time)
        go crawler.Crawl(ctx)
        time.Sleep(2 * time.Second)

        // Check Redis cache
        body, err := crawler.redisClient.Get(ctx, "cache:https://example.com").Result()
        assert.NoError(t, err)
        assert.Contains(t, body, "Example Domain")

        // Check seen URLs
        seen, err := crawler.redisClient.SIsMember(ctx, "seen_urls", "https://example.com").Result()
        assert.NoError(t, err)
        assert.True(t, seen)
    })
}