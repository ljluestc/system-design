package api

import (
    "log"
    "net/http"

    "github.com/gorilla/mux"
)

// StartServer starts the HTTP server
func StartServer(port string, suggestionHandler, queryHandler http.HandlerFunc) {
    r := mux.NewRouter()
    r.HandleFunc("/suggestions", suggestionHandler).Methods("GET")

    log.Printf("Starting server on :%s", port)
    if err := http.ListenAndServe(":"+port, r); err != nil {
        log.Fatalf("Server failed: %v", err)
    }
}