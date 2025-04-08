package controllers

import (
    "encoding/json"
    "net/http"
    "database-demo/services"
)

func MonitoringStatsHandler(w http.ResponseWriter, r *http.Request) {
    stats := services.GetMonitoringStats()
    json.NewEncoder(w).Encode(stats)
}