package api

import (
    "encoding/json"
    "net/http"
    "key-value-store/internal/db"
)

func GetHandler(w http.ResponseWriter, r *http.Request) {
    key := r.URL.Query().Get("key")
    if key == "" {
        http.Error(w, "Missing key", http.StatusBadRequest)
        return
    }
    value, err := db.Get(key)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(map[string]string{"value": value})
}

func PutHandler(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Key   string `json:"key"`
        Value string `json:"value"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }
    if err := db.Put(req.Key, req.Value); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusCreated)
}