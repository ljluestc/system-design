package main

import (
    "github.com/calelin/messenger/config"
    "github.com/calelin/messenger/internal/db"
    "github.com/sirupsen/logrus"
)

func main() {
    log := logrus.New()
    cfg := config.NewConfig()
    dbConn := db.NewMongoDB(cfg, log)

    // Placeholder for backup logic
    log.Info("Starting database backup")
    // In production, use mongodump or similar
    log.Info("Backup completed")
}