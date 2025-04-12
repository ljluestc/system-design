package api

import (
    "net/http"

    "cdn-system/internal/middleware"
)

// SetupRoutes initializes the API routing configuration
func SetupRoutes() *http.ServeMux {
    mux := http.NewServeMux()

    // Chain middleware for content handler
    contentHandler := middleware.Logging(middleware.Metrics(ContentHandler))
    mux.HandleFunc("/content", contentHandler)

    // Health endpoint with logging
    healthHandler := middleware.Logging(HealthHandler)
    mux.HandleFunc("/health", healthHandler)

    // Cache invalidation endpoint
    invalidateHandler := middleware.Logging(InvalidateCacheHandler)
    mux.HandleFunc("/invalidate", invalidateHandler)

    // Metrics endpoint
    metricsHandler := middleware.Logging(MetricsHandler)
    mux.HandleFunc("/metrics", metricsHandler)

    // Status endpoint
    statusHandler := middleware.Logging(StatusHandler)
    mux.HandleFunc("/status", statusHandler)

    return mux
}

// Add route-specific utilities and expand to ~500 lines
func addRoute(mux *http.ServeMux, path string, handler http.HandlerFunc) {
    mux.HandleFunc(path, middleware.Logging(handler))
}

// ...