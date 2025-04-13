package db

import (
    "github.com/calelin/messenger/internal/encoder"
    "github.com/sirupsen/logrus"
    "time"
)

// Sequencer generates unique IDs
type Sequencer struct {
    db  *MongoDB
    log *logrus.Logger
}

// NewSequencer creates a new Sequencer
func NewSequencer(db *MongoDB, log *logrus.Logger) *Sequencer {
    return &Sequencer{db: db, log: log}
}

// GenerateShortURL creates a short URL
func (s *Sequencer) GenerateShortURL(originalURL, customAlias string, expiryDate *time.Time) (string, error) {
    var shortKey string
    if customAlias != "" {
        // Validate and use custom alias
        if err := s.db.ReserveCustomAlias(customAlias, originalURL, expiryDate); err != nil {
            return "", err
        }
        shortKey = customAlias
    } else {
        // Generate unique ID
        id, err := s.db.GetNextID()
        if err != nil {
            return "", err
        }
        // Encode to base-58
        shortKey = encoder.ToBase58(id)
        if err := s.db.SaveURL(shortKey, originalURL, expiryDate); err != nil {
            return "", err
        }
    }
    return "http://tiny.url/" + shortKey, nil
}