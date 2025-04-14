package main

import (
    "github.com/calelin/instagram/config"
    "github.com/calelin/instagram/internal/api"
    "github.com/calelin/instagram/internal/auth"
    "github.com/calelin/instagram/internal/cache"
    "github.com/calelin/instagram/internal/db"
    "github.com/calelin/instagram/internal/queue"
    "github.com/calelin/instagram/internal/storage"
    "github.com/calelin/instagram/internal/utils"
    "github.com/calelin/instagram/pkg/server"
)

func main() {
    log := utils.NewLogger()
    cfg := config.NewConfig()
    postgresDB := db.NewPostgresDB(cfg, log)
    cassandraDB := db.NewCassandraDB(cfg, log)
    redisClient := cache.NewRedisClient(cfg, log)
    s3Storage := storage.NewS3Storage(cfg, log)
    kafkaProducer := queue.NewKafkaProducer(cfg, log)
    authService := auth.NewAuthService(postgresDB, log)

    // Initialize handler
    h := api.NewHandler(authService, postgresDB, cassandraDB, redisClient, s3Storage, kafkaProducer, log)

    // Start server
    server.Start(cfg, h, log)
}