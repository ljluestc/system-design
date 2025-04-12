package edge

import (
    "log"
    "time"
)

// UpdateHealth updates the health timestamp
func (s *EdgeServer) UpdateHealth() {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.LastHealth = time.Now()
    log.Printf("Health updated for %s", s.ID)
}

// IsHealthy checks if the server is healthy
func (s *EdgeServer) IsHealthy() bool {
    s.mu.Lock()
    defer s.mu.Unlock()
    return time.Since(s.LastHealth) < 2*time.Minute
}

// Expand with detailed health check logic
// ...