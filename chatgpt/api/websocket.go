package api

import "fmt"

type WebSocketServer struct{}

func NewWebSocketServer() *WebSocketServer {
    return &WebSocketServer{}
}

func (ws *WebSocketServer) Stream(userID, prompt string) {
    fmt.Printf("Streaming to %s: %s\n", userID, prompt)
}