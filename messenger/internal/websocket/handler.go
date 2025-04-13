package websocket

import (
    "log"
    "net/http"
    "github.com/gorilla/websocket"
    "messenger/internal/auth"
    "messenger/internal/chat"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool { return true },
}

type Handler struct {
    chatSvc  *chat.Service
    authSvc  *auth.Service
    conns    map[string]*websocket.Conn
}

func NewHandler(chatSvc *chat.Service, authSvc *auth.Service) *Handler {
    return &Handler{
        chatSvc: chatSvc,
        authSvc: authSvc,
        conns:   make(map[string]*websocket.Conn),
    }
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
    // Upgrade HTTP to WebSocket
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Printf("WebSocket upgrade failed: %v", err)
        return
    }
    defer conn.Close()

    // Get user ID from query param (simplified auth)
    userID := r.URL.Query().Get("user_id")
    if userID == "" {
        log.Println("Missing user_id")
        return
    }
    h.conns[userID] = conn

    // Handle incoming messages
    for {
        var msg chat.Message
        if err := conn.ReadJSON(&msg); err != nil {
            log.Printf("Read error: %v", err)
            break
        }
        msg.SenderID = userID
        if _, err := h.chatSvc.SendMessage(msg.SenderID, msg.RecipientID, msg.GroupID, msg.Content); err != nil {
            log.Printf("Send error: %v", err)
            continue
        }
        // Send to recipient if online
        if conn, ok := h.conns[msg.RecipientID]; ok {
            conn.WriteJSON(msg)
        }
    }
    delete(h.conns, userID)
}