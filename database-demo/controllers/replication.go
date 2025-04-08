package controllers

import (
    "encoding/json"
    "net/http"
    "database-demo/services"
)

func ReplicationSyncHandler(w http.ResponseWriter, r *http.Request) {
    if err := services.SyncReplicas(); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(map[string]string{"message": "Replicas synced"})
}