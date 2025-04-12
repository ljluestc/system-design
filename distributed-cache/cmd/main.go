package main

import (
    "distributed-cache/internal/cache"
    "distributed-cache/internal/server"
    "distributed-cache/internal/utils"
    "log"
    "os"
    "os/signal"
    "syscall"
)

func main() {
    // Load configuration
    config := utils.LoadConfig()
    if config == nil {
        log.Fatal("Failed to load configuration")
    }

    // Initialize hash ring with 3 replicas per node
    hashRing := cache.NewHashRing(3)
    nodes := make(map[string]*cache.CacheNode)

    // Add nodes to the hash ring and create cache instances
    for _, nodeAddr := range config.Nodes {
        hashRing.AddNode(nodeAddr)
        nodes[nodeAddr] = cache.NewCacheNode(1000) // Capacity of 1000 items
    }

    // Start the server in a goroutine
    go func() {
        log.Printf("Starting server on port %s", config.ServerPort)
        if err := server.StartServer(config.ServerPort, hashRing, nodes); err != nil {
            log.Fatalf("Server failed: %v", err)
        }
    }()

    // Handle graceful shutdown
    sigs := make(chan os.Signal, 1)
    signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
    <-sigs
    log.Println("Shutting down gracefully...")
}