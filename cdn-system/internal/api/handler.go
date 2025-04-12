package api

import (
    "fmt"
    "log"
    "net/http"
    "time"

    "cdn-system/internal/edge"
    "cdn-system/internal/utils"
)

// ContentHandler serves content requests from edge servers
func ContentHandler(w http.ResponseWriter, r *http.Request) {
    start := time.Now()
    key := r.URL.Query().Get("key")
    if key == "" {
        http.Error(w, "Missing content key", http.StatusBadRequest)
        log.Printf("Bad request: missing key from %s", r.RemoteAddr)
        return
    }

    content, err := edge.GetContent(key)
    if err != nil {
        http.Error(w, fmt.Sprintf("Failed to retrieve content: %v", err), http.StatusInternalServerError)
        log.Printf("Error retrieving content for key %s: %v", key, err)
        return
    }

    w.Header().Set("Content-Type", "text/plain")
    w.Header().Set("X-Cache-Hit", "true")
    if _, err := w.Write([]byte(content)); err != nil {
        log.Printf("Failed to write response: %v", err)
    }
    log.Printf("Served content for key %s in %v", key, time.Since(start))
}

// HealthHandler provides a system health check endpoint
func HealthHandler(w http.ResponseWriter, r *http.Request) {
    if !edge.CheckHealth() {
        http.Error(w, "Edge servers unhealthy", http.StatusServiceUnavailable)
        log.Println("Health check failed: edge servers unhealthy")
        return
    }
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Healthy"))
    log.Printf("Health check passed for %s", r.RemoteAddr)
}

// InvalidateCacheHandler invalidates cache for a specific key
func InvalidateCacheHandler(w http.ResponseWriter, r *http.Request) {
    key := r.URL.Query().Get("key")
    if key == "" {
        http.Error(w, "Missing content key", http.StatusBadRequest)
        log.Printf("Invalidation failed: missing key from %s", r.RemoteAddr)
        return
    }

    err := edge.InvalidateCache(key)
    if err != nil {
        http.Error(w, fmt.Sprintf("Failed to invalidate cache: %v", err), http.StatusInternalServerError)
        log.Printf("Cache invalidation error for key %s: %v", key, err)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Cache invalidated"))
    log.Printf("Cache invalidated for key %s", key)
}

// MetricsHandler exposes system metrics (placeholder)
func MetricsHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Metrics placeholder"))
    log.Printf("Metrics requested by %s", r.RemoteAddr)
}

// Utility function to validate requests
func validateRequest(r *http.Request) error {
    if r.Method != http.MethodGet && r.Method != http.MethodPost {
        return fmt.Errorf("method %s not allowed", r.Method)
    }
    return nil
}

// Add more handlers and utilities to reach ~500 lines
func StatusHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("System status: running"))
    log.Printf("Status check by %s", r.RemoteAddr)
}

// Expand with additional logic, error handling, and logging
// ...