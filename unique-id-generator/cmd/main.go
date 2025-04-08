package main

import (
    "fmt"
    "log"
    "time"

    "unique-id-generator/internal/config"
    "unique-id-generator/internal/sequencer"
)

func main() {
    // Load configuration
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatalf("Failed to load configuration: %v", err)
    }

    // Create a new Sequencer
    seq, err := sequencer.NewSequencer(cfg.NodeID, cfg.Epoch)
    if err != nil {
        log.Fatalf("Failed to create sequencer: %v", err)
    }

    // Generate and print 10 unique IDs
    for i := 0; i < 10; i++ {
        id, err := seq.GenerateID()
        if err != nil {
            log.Fatalf("Error generating ID: %v", err)
        }
        fmt.Printf("Generated ID: %d\n", id)
        time.Sleep(100 * time.Millisecond) // Simulate time passing
    }
}