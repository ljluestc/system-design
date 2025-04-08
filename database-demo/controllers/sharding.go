package controllers

import (
    "encoding/json"
    "net/http"
    "database-demo/services"
)

type InsertRequest struct {
    UserID string `json:"userId"`
    Data   string `json:"data"`
}

func ShardingInsertHandler(w http.ResponseWriter, r *http.Request) {
    var req InsertRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }
    shardID := services.GetShardID(req.UserID)
    if err := services.InsertToShard(shardID, req.UserID, req.Data); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{"message": "Data inserted", "shardId": string(shardID)})
}

func ShardingGetHandler(w http.ResponseWriter, r *http.Request) {
    userID := r.URL.Query().Get("userId")
    if userID == "" {
        http.Error(w, "Missing userId", http.StatusBadRequest)
        return
    }
    shardID := services.GetShardID(userID)
    data, err := services.QueryShard(shardID, userID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(data)
}