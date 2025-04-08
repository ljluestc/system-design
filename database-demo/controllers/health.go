package controllers

import (
    "encoding/json"
    "net/http"
    "database-demo/services"
)

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
    if services.IsHealthy() {
        json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
    } else {
        http.Error(w, "Unhealthy", http.StatusServiceUnavailable)
    }
}