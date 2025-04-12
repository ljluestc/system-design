package api

import (
    "net/http"
)

func SetupRoutes() *http.ServeMux {
    mux := http.NewServeMux()
    mux.HandleFunc("/get", GetHandler)
    mux.HandleFunc("/put", PutHandler)
    return mux
}