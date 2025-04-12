package middleware

import (
    "log"
    "net/http"
    "time"
)

// Logging adds request logging
func Logging(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        log.Printf("%s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
        next(w, r)
        log.Printf("Completed %s %s in %v", r.Method, r.URL.Path, time.Since(start))
    }
}

// Expand with detailed logging
// ...