package middleware

import "net/http"

func RateLimit(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Placeholder for rate limiting
        next(w, r)
    }
}