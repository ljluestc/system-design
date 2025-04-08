package sequencer

import (
    "fmt"
    "sync"
    "time"

    "unique-id-generator/internal/utils"
)

// Sequencer generates time-sortable unique IDs
type Sequencer struct {
    nodeID        int64
    epoch         int64
    sequence      int64
    lastTimestamp int64
    maxSequence   int64
    mu            sync.Mutex
}

// NewSequencer creates a new Sequencer instance
func NewSequencer(nodeID, epoch int64) (*Sequencer, error) {
    if nodeID < 0 || nodeID > 1023 {
        return nil, fmt.Errorf("nodeID must be between 0 and 1023")
    }
    return &Sequencer{
        nodeID:        nodeID,
        epoch:         epoch,
        sequence:      0,
        lastTimestamp: 0,
        maxSequence:   4096, // 12 bits for sequence number
    }, nil
}

// GenerateID generates a new unique ID
func (s *Sequencer) GenerateID() (int64, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    timestamp := utils.GetCurrentTimestamp()

    if timestamp < s.lastTimestamp {
        return 0, fmt.Errorf("clock moved backwards: %d < %d", timestamp, s.lastTimestamp)
    }

    if timestamp == s.lastTimestamp {
        s.sequence = (s.sequence + 1) % s.maxSequence
        if s.sequence == 0 {
            // Wait for the next millisecond if sequence overflows
            for timestamp <= s.lastTimestamp {
                timestamp = utils.GetCurrentTimestamp()
                time.Sleep(1 * time.Millisecond)
            }
        }
    } else {
        s.sequence = 0
    }

    s.lastTimestamp = timestamp

    // Compose the 64-bit ID: 41-bit timestamp | 10-bit nodeID | 12-bit sequence
    id := ((timestamp - s.epoch) << 22) | (s.nodeID << 12) | s.sequence
    return id, nil
}