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

    // Seed sample URLs
    urls := []struct {
        shortKey    string
        originalURL string
    }{
        {"abc123", "https://example.com"},
        {"xyz789", "https://google.com"},
    }
    for _, url := range urls {
        if err := dbConn.SaveURL(url.shortKey, url.originalURL, nil); err != nil {
            log.Errorf("Failed to seed %s: %v", url.shortKey, err)
        } else {
            log.Infof("Seeded %s", url.shortKey)
        }
    }
}