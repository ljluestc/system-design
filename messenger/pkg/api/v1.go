package api

import (
    "github.com/gorilla/mux"
    "messenger/internal/auth"
    "messenger/internal/chat"
    "messenger/internal/media"
    "messenger/internal/notification"
    "messenger/internal/websocket"
)

func NewRouter(authSvc *auth.Service, chatSvc *chat.Service, mediaSvc *media.Service, notifSvc *notification.Service, wsHandler *websocket.Handler) *mux.Router {
    r := mux.NewRouter()
    r.HandleFunc("/ws", wsHandler.Handle)
    r.HandleFunc("/login", authSvc.LoginHandler).Methods("POST")
    r.HandleFunc("/users", authSvc.CreateUserHandler).Methods("POST")
    r.HandleFunc("/messages", chatSvc.SendMessageHandler).Methods("POST")
    r.HandleFunc("/media/upload", mediaSvc.UploadHandler).Methods("POST")
    return r
}