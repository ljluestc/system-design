package server

import (
    "github.com/calelin/messenger/config"
    "github.com/calelin/messenger/internal/api"
    "github.com/calelin/messenger/internal/auth"
    "github.com/calelin/messenger/internal/cache"
    "github.com/calelin/messenger/internal/ratelimit"
    "github.com/gorilla/mux"
    "github.com/sirupsen/logrus"
    "net/http"
)

// Start starts the HTTP server
func Start(cfg *config.Config, authService *auth.AuthService, sequencer *db.Sequencer, cacheClient *cache.MemcachedClient, rateLimiter *ratelimit.Limiter, log *logrus.Logger) {
    r := mux.NewRouter()
    h := api.NewHandler(sequencer, cacheClient, log)
    m := api.NewMiddleware(rateLimiter, authService, log)

    // Routes with middleware
    r.HandleFunc("/shorten", m.RateLimitMiddleware(m.AuthMiddleware(h.ShortenHandler))).Methods("POST")
    r.HandleFunc("/{shortKey}", h.RedirectHandler).Methods("GET")
    r.HandleFunc("/delete/{shortKey}", m.RateLimitMiddleware(m.AuthMiddleware(h.DeleteHandler))).Methods("DELETE")

    log.Printf("Server running on :%s", cfg.Port)
    http.ListenAndServe(":"+cfg.Port, r)
}