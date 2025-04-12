package server

import (
    "net/http"
    "github.com/gorilla/mux"
    "distributed-cache/internal/cache"
    "distributed-cache/pkg/api"
)

func StartServer(port string, hashRing *cache.HashRing, nodes map[string]*cache.CacheNode) error {
    router := mux.NewRouter()
    router.Use(api.LoggingMiddleware) // Add logging middleware
    router.HandleFunc("/cache/{key}", api.GetCacheHandler(hashRing, nodes)).Methods("GET")
    router.HandleFunc("/cache/{key}", api.SetCacheHandler(hashRing, nodes)).Methods("PUT")
    return http.ListenAndServe(":"+port, router)
}