package config

import (
    "errors"
    "os"
    "strconv"
)

// Config holds the configuration for the sequencer
type Config struct {
    NodeID int64 // Unique identifier for the node
    Epoch  int64 // Custom epoch timestamp in milliseconds
}

// LoadConfig loads the configuration from environment variables
func LoadConfig() (*Config, error) {
    // Get NODE_ID from environment
    nodeIDStr := os.Getenv("NODE_ID")
    if nodeIDStr == "" {
        return nil, errors.New("NODE_ID environment variable is required")
    }
    nodeID, err := strconv.ParseInt(nodeIDStr, 10, 64)
    if err != nil {
        return nil, errors.New("invalid NODE_ID, must be an integer")
    }

    // Get EPOCH from environment or use default (January 1, 2020)
    epochStr := os.Getenv("EPOCH")
    if epochStr == "" {
        epochStr = "1577836800000" // Default epoch: 2020-01-01 00:00:00 UTC
    }
    epoch, err := strconv.ParseInt(epochStr, 10, 64)
    if err != nil {
        return nil, errors.New("invalid EPOCH, must be a timestamp in milliseconds")
    }

    return &Config{
        NodeID: nodeID,
        Epoch:  epoch,
    }, nil
}