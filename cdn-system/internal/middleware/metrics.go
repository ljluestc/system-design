package middleware

import (
    "log"
    "net/http"
    "time"
)

// Metrics adds metrics collection
func Metrics(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next(w, r)
        log.Printf("Metrics: %s %s took %v", r.Method, r.URL.Path, time.Since(start))
    }
}

// Expand with Prometheus integration
// ...