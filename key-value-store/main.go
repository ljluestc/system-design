// main.go
package main

import (
    "key-value-store/api"
    "key-value-store/config"
    "key-value-store/logger"
    "key-value-store/node"
    "key-value-store/shard"
    "log"
)

func main() {
    logger.SetupLogger()
    log.Println("Starting key-value store cluster...")

    cfg := config.LoadConfig()
    nodes := make([]*node.Node, 0)

    // Initialize shards with primaries and replicas
    for i := 0; i < cfg.NumShards; i++ {
        primary := node.NewNode(i, true, nil)
        nodes = append(nodes, primary)
        for j := 0; j < cfg.ReplicasPerShard; j++ {
            replica := node.NewNode(i, false, primary)
            nodes = append(nodes, replica)
            primary.AddReplica(replica)
        }
    }

    router := shard.NewRouter(nodes)
    server := api.NewServer(router)

    log.Printf("Cluster started: %d shards, %d replicas each", cfg.NumShards, cfg.ReplicasPerShard)

    // Start API server
    if err := server.Start(); err != nil {
        log.Fatalf("API server failed: %v", err)
    }
}