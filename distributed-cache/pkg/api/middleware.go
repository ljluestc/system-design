package api

import (
    "net/http"
    "distributed-cache/internal/utils"
)

func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        utils.Logger.Printf("%s %s", r.Method, r.URL.Path)
        next.ServeHTTP(w, r)
    })
}