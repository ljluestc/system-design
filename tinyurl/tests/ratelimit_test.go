package tests

import (
    "github.com/calelin/messenger/config"
    "github.com/calelin/messenger/internal/ratelimit"
    "github.com/sirupsen/logrus"
    "testing"
)

func TestRateLimiter(t *testing.T) {
    log := logrus.New()
    cfg := config.NewConfig()
    limiter := ratelimit.NewLimiter(cfg, log)
    for i := 0; i < 100; i++ {
        if !limiter.Allow("test-key") {
            t.Error("Expected request to be allowed")
        }
    }
    if limiter.Allow("test-key") {
        t.Error("Expected request to be blocked")
    }
}