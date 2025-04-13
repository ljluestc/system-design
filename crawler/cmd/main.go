package main

import (
    "log"
    "net/http"
    "crawler/config"
    "crawler/internal/dedup"
    "crawler/internal/dns"
    "crawler/internal/extractor"
    "crawler/internal/fetcher"
    "crawler/internal/scheduler"
    "crawler/internal/storage"
    "crawler/internal/worker"
    "crawler/pkg/api"
    "crawler/pkg/db"
)

func main() {
    // Load configuration
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    // Initialize Postgres for URL storage
    dbConn, err := db.NewPostgres(cfg.PostgresURL)
    if err != nil {
        log.Fatalf("Failed to connect to Postgres: %v", err)
    }
    defer dbConn.Close()

    // Initialize Redis for caching
    redisConn, err := storage.NewRedis(cfg.RedisURL)
    if err != nil {
        log.Fatalf("Failed to connect to Redis: %v", err)
    }

    // Initialize services
    schedulerSvc := scheduler.New(dbConn, redisConn)
    dnsSvc := dns.New(cfg.DNSCacheTTL)
    fetcherSvc := fetcher.New(cfg.ProxyList)
    extractorSvc := extractor.New()
    dedupSvc := dedup.New(dbConn)
    storageSvc := storage.New(cfg.S3Bucket)

    // Start worker manager
    manager := worker.NewManager(schedulerSvc, dnsSvc, fetcherSvc, extractorSvc, dedupSvc, storageSvc, cfg.WorkerCount)
    go manager.Start()

    // Start API server
    router := api.NewRouter(schedulerSvc)
    log.Printf("Starting API server on :%s", cfg.APIPort)
    if err := http.ListenAndServe(":"+cfg.APIPort, router); err != nil {
        log.Fatalf("API server failed: %v", err)
    }
}