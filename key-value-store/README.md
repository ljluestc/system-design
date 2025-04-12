package main

import (
    "log"
    "net/http"
    "key-value-store/internal/api"
    "key-value-store/internal/config"
    "key-value-store/internal/db"
    "key-value-store/internal/monitoring"
)

func main() {
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }
    if err := db.InitNodes(cfg); err != nil {
        log.Fatalf("Failed to initialize nodes: %v", err)
    }
    monitoring.Setup()
    router := api.SetupRoutes()
    log.Printf("Starting server on :%s", cfg.Port)
    if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
        log.Fatalf("Server failed: %v", err)
    }
}