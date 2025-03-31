package main

import (
    "encoding/json"
    "log"
    "net/http"
    "sync"

    "github.com/gorilla/websocket"
)

// Document represents a document with versioning
type Document struct {
    ID       string
    Content  string
    Version  int
    Mutex    sync.Mutex
}

// Operation represents an edit operation
type Operation struct {
    Type     string // "insert" or "delete"
    Position int
    Char     rune // For insert
    ClientID string
    Version  int
}

// Server manages documents and clients
type Server struct {
    documents map[string]*Document
    clients   map[*websocket.Conn]string
    mu        sync.Mutex
}

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool { return true },
}

func NewServer() *Server {
    return &Server{
        documents: make(map[string]*Document),
        clients:   make(map[*websocket.Conn]string),
    }
}

func (s *Server) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println("WebSocket upgrade error:", err)
        return
    }
    defer conn.Close()

    // Read document ID
    var msg struct {
        DocID string `json:"doc_id"`
    }
    if err := conn.ReadJSON(&msg); err != nil {
        log.Println("Error reading doc_id:", err)
        return
    }

    // Register client
    s.mu.Lock()
    s.clients[conn] = msg.DocID
    s.mu.Unlock()

    // Send initial document state
    doc, ok := s.documents[msg.DocID]
    if !ok {
        doc = &Document{ID: msg.DocID, Content: "", Version: 0}
        s.documents[msg.DocID] = doc
    }
    conn.WriteJSON(struct {
        Type    string `json:"type"`
        Content string `json:"content"`
        Version int    `json:"version"`
    }{"init", doc.Content, doc.Version})

    // Handle operations
    for {
        var op Operation
        if err := conn.ReadJSON(&op); err != nil {
            log.Println("Error reading operation:", err)
            break
        }
        s.applyOperation(msg.DocID, op, conn)
    }
}

func (s *Server) applyOperation(docID string, op Operation, sender *websocket.Conn) {
    s.mu.Lock()
    doc, ok := s.documents[docID]
    if !ok {
        s.mu.Unlock()
        return
    }
    doc.Mutex.Lock()
    defer doc.Mutex.Unlock()

    // Basic conflict check
    if op.Version < doc.Version {
        log.Println("Rejecting outdated operation")
        return
    }

    // Apply operation
    switch op.Type {
    case "insert":
        if op.Position <= len(doc.Content) {
            doc.Content = doc.Content[:op.Position] + string(op.Char) + doc.Content[op.Position:]
        }
    case "delete":
        if op.Position < len(doc.Content) {
            doc.Content = doc.Content[:op.Position] + doc.Content[op.Position+1:]
        }
    }
    doc.Version++
    s.mu.Unlock()

    // Broadcast to other clients
    s.mu.Lock()
    for client, id := range s.clients {
        if id == docID && client != sender {
            client.WriteJSON(op)
        }
    }
    s.mu.Unlock()
}

func main() {
    server := NewServer()
    http.HandleFunc("/ws", server.HandleWebSocket)
    log.Println("Server starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}