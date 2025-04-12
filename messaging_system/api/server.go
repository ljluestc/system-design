package api

import (
    "net/http"
    "messaging_system/services"

    "github.com/gorilla/mux"
    "github.com/sirupsen/logrus"
)

func StartServer(port string, msgService *services.MessagingService, log *logrus.Logger) {
    r := mux.NewRouter()
    // Add routes here
    log.Printf("Server running on :%s", port)
    http.ListenAndServe(":"+port, r)
}