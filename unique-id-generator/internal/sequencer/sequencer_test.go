package sequencer_test

import (
    "testing"
    "time"

    "unique-id-generator/internal/sequencer"
)

func TestSequencer_GenerateID(t *testing.T) {
    seq, err := sequencer.NewSequencer(1, time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).UnixNano()/1e6)
    if err != nil {
        t.Fatalf("Failed to create sequencer: %v", err)
    }

    ids := make([]int64, 10)
    for i := 0; i < 10; i++ {
        id, err := seq.GenerateID()
        if err != nil {
            t.Fatalf("Error generating ID: %v", err)
        }
        ids[i] = id
        time.Sleep(100 * time.Millisecond)
    }

    // Check if IDs are unique and increasing
    for i := 1; i < len(ids); i++ {
        if ids[i] <= ids[i-1] {
            t.Errorf("IDs are not time-sortable: %d <= %d", ids[i], ids[i-1])
        }
    }
}

func TestSequencer_ClockBackwards(t *testing.T) {
    seq, _ := sequencer.NewSequencer(1, time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).UnixNano()/1e6)
    // Simulate clock moving backwards
    seq.lastTimestamp = time.Now().UnixNano()/1e6 + 1000 // Set last timestamp to future
    _, err := seq.GenerateID()
    if err == nil {
        t.Error("Expected error for clock moving backwards, got none")
    }
}