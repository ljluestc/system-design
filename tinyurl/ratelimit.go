package main

import (
    "sync"
    "time"
)

// RateLimiter implements a token bucket rate limiter
type RateLimiter struct {
    rate       int           // tokens per interval
    interval   time.Duration // refill interval
    capacity   int           // max tokens
    tokens     int           // current tokens
    lastRefill time.Time     // last refill time
    mu         sync.Mutex
}

// NewRateLimiter initializes a rate limiter
func NewRateLimiter(rate int, interval time.Duration) *RateLimiter {
    return &RateLimiter{
        rate:       rate,
        interval:   interval,
        capacity:   rate,
        tokens:     rate,
        lastRefill: time.Now(),
    }
}

// Allow checks if a request is allowed
func (r *RateLimiter) Allow() bool {
    r.mu.Lock()
    defer r.mu.Unlock()

    now := time.Now()
    r.refill(now)

    if r.tokens > 0 {
        r.tokens--
        return true
    }
    return false
}

// refill adds tokens based on elapsed time
func (r *RateLimiter) refill(now time.Time) {
    elapsed := now.Sub(r.lastRefill)
    tokensToAdd := int(elapsed / r.interval * time.Duration(r.rate))
    if tokensToAdd > 0 {
        r.tokens = min(r.capacity, r.tokens+tokensToAdd)
        r.lastRefill = now
    }
}

// min returns the minimum of two integers
func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}