package middleware

import "net/http"

func Metrics(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Placeholder for metrics logic
        next(w, r)
    }
}