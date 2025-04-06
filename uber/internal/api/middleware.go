package api

import (
    "net/http"
    "uber/internal/auth"
    "uber/internal/utils"
)

// LoggingMiddleware logs incoming requests
func LoggingMiddleware(logger *utils.Logger) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            logger.Info("Request received:", r.Method, r.URL.Path)
            next.ServeHTTP(w, r)
        })
    }
}

// AuthMiddleware validates JWT tokens and sets user ID in the request header
func AuthMiddleware(authSvc *auth.AuthService, logger *utils.Logger) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            token := r.Header.Get("Authorization")
            if token == "" {
                logger.Error("No token provided")
                http.Error(w, "Unauthorized", http.StatusUnauthorized)
                return
            }

            userID, err := authSvc.ValidateToken(token)
            if err != nil {
                logger.Error("Invalid token:", err)
                http.Error(w, "Unauthorized", http.StatusUnauthorized)
                return
            }

            r.Header.Set("UserID", userID)
            next.ServeHTTP(w, r)
        })
    }
}

// RateLimitMiddleware provides basic rate limiting (placeholder)
func RateLimitMiddleware(logger *utils.Logger) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Placeholder for rate limiting logic (e.g., token bucket)
            logger.Info("Rate limit check passed")
            next.ServeHTTP(w, r)
        })
    }
}