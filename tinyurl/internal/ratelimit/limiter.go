package ratelimit

import (
    "github.com/calelin/messenger/config"
    "github.com/sirupsen/logrus"
    "sync"
    "time"
)

// Limiter enforces rate limits
type Limiter struct {
    limits  map[string]int
    counters map[string]int
    mu       sync.Mutex
    log      *logrus.Logger
    window   time.Duration
}

// NewLimiter creates a new Limiter
func NewLimiter(cfg *config.Config, log *logrus.Logger) *Limiter {
    l := &Limiter{
        limits:   map[string]int{},
        counters: map[string]int{},
        log:      log,
        window:   time.Second * 60, // 1 minute window
    }
    go l.resetCounters()
    return l
}

// Allow checks if a request is allowed
func (l *Limiter) Allow(apiKey string) bool {
    l.mu.Lock()
    defer l.mu.Unlock()

    // Default limit: 100 requests per minute
    limit := 100
    if customLimit, exists := l.limits[apiKey]; exists {
        limit = customLimit
    }

    count, exists := l.counters[apiKey]
    if !exists {
        l.counters[apiKey] = 1
        return true
    }

    if count >= limit {
        l.log.Warnf("Rate limit exceeded for %s", apiKey)
        return false
    }

    l.counters[apiKey]++
    return true
}

// resetCounters clears counters periodically
func (l *Limiter) resetCounters() {
    ticker := time.NewTicker(l.window)
    for range ticker.C {
        l.mu.Lock()
        l.counters = make(map[string]int)
        l.mu.Unlock()
    }
}