package main

import (
    "chatgpt/api"
    "chatgpt/config"
    "chatgpt/inference"
    "chatgpt/logger"
    "chatgpt/monitoring"
    "log"
)

func main() {
    logger.SetupLogger()
    log.Println("Starting ChatGPT service...")

    cfg := config.LoadConfig()
    infService := inference.NewService(cfg.NumShards)
    router := inference.NewRouter(infService)
    server := api.NewServer(router)

    go monitoring.StartMetrics()

    log.Printf("Service started with %d shards", cfg.NumShards)
    if err := server.Start(); err != nil {
        log.Fatalf("Server failed: %v", err)
    }
}