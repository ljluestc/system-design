package main

import (
    "github.com/calelin/tinyurl/config"
    "github.com/calelin/tinyurl/internal/api"
    "github.com/calelin/tinyurl/internal/auth"
    "github.com/calelin/tinyurl/internal/cache"
    "github.com/calelin/tinyurl/internal/db"
    "github.com/calelin/tinyurl/internal/ratelimit"
    "github.com/calelin/tinyurl/internal/utils"
    "github.com/calelin/tinyurl/pkg/server"
)

func main() {
    log := utils.NewLogger()
    cfg := config.NewConfig()
    mongoDB := db.NewMongoDB(cfg)
    cacheClient := cache.NewMemcachedClient(cfg)
    sequencer := db.NewSequencer(mongoDB)
    rateLimiter := ratelimit.NewLimiter(cfg, log)
    authService := auth.NewAuthService(mongoDB, log)

    // Start the server
    server.Start(cfg, authService, sequencer, cacheClient, rateLimiter, log)
}