package api

import (
    "io/ioutil"
    "net/http"
    "time"
    "github.com/gorilla/mux"
    "distributed-cache/internal/cache"
)

func GetCacheHandler(hashRing *cache.HashRing, nodes map[string]*cache.CacheNode) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        key := vars["key"]
        nodeAddr := hashRing.GetNode(key)
        node := nodes[nodeAddr]
        if val, ok := node.Get(key); ok {
            w.Write([]byte(val))
        } else {
            w.WriteHeader(http.StatusNotFound)
            w.Write([]byte("Key not found"))
        }
    }
}

func SetCacheHandler(hashRing *cache.HashRing, nodes map[string]*cache.CacheNode) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        key := vars["key"]
        value, err := ioutil.ReadAll(r.Body)
        if err != nil {
            http.Error(w, "Invalid request body", http.StatusBadRequest)
            return
        }
        nodeAddr := hashRing.GetNode(key)
        node := nodes[nodeAddr]
        node.Set(key, string(value), 3600*time.Second) // 1-hour TTL
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("OK"))
    }
}