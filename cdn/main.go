package main

import (
    "cdn/api"
    "cdn/config"
    "cdn/edge"
    "cdn/monitoring"
    "cdn/routing"
    "cdn/scrubber"
    "log"
)

func main() {
    log.Println("Starting CDN...")
    cfg := config.LoadConfig()
    edges := make([]*edge.Server, cfg.NumEdges)
    for i := 0; i < cfg.NumEdges; i++ {
        edges[i] = edge.NewServer(i, cfg.CacheCapacity)
    }
    scrubberService := scrubber.NewService()
    router := routing.NewRouter(edges)
    server := api.NewServer(router, scrubberService)
    go monitoring.StartMetrics()
    if err := server.Start(); err != nil {
        log.Fatalf("Server failed: %v", err)
    }
}