package api

import (
    "github.com/calelin/messenger/internal/auth"
    "github.com/calelin/messenger/internal/ratelimit"
    "github.com/sirupsen/logrus"
    "net/http"
)

// Middleware handles rate limiting and authentication
type Middleware struct {
    rateLimiter *ratelimit.Limiter
    authService *auth.AuthService
    log         *logrus.Logger
}

// NewMiddleware creates a new Middleware
func NewMiddleware(rateLimiter *ratelimit.Limiter, authService *auth.AuthService, log *logrus.Logger) *Middleware {
    return &Middleware{
        rateLimiter: rateLimiter,
        authService: authService,
        log:         log,
    }
}

// RateLimitMiddleware limits requests per API key
func (m *Middleware) RateLimitMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        apiKey := r.Header.Get("X-API-Key")
        if apiKey == "" {
            m.log.Warn("Missing API key")
            respondWithError(w, http.StatusUnauthorized, "Missing API key")
            return
        }

        if !m.rateLimiter.Allow(apiKey) {
            m.log.Warnf("Rate limit exceeded for API key: %s", apiKey)
            respondWithError(w, http.StatusTooManyRequests, "Rate limit exceeded")
            return
        }

        next(w, r)
    }
}

// AuthMiddleware validates API keys
func (m *Middleware) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if err := m.authService.ValidateAPIKey(r.Header.Get("X-API-Key")); err != nil {
            m.log.Errorf("Invalid API key: %v", err)
            respondWithError(w, http.StatusUnauthorized, "Invalid API key")
            return
        }
        next(w, r)
    }
}