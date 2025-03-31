package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
    "strings"

    "github.com/gorilla/websocket"
)

type Client struct {
    conn  *websocket.Conn
    docID string
}

func NewClient(serverURL, docID string) (*Client, error) {
    conn, _, err := websocket.DefaultDialer.Dial(serverURL, nil)
    if err != nil {
        return nil, err
    }
    return &Client{conn: conn, docID: docID}, nil
}

func (c *Client) Start() {
    // Send docID
    c.conn.WriteJSON(struct {
        DocID string `json:"doc_id"`
    }{c.docID})

    // Handle server messages
    go func() {
        for {
            var msg map[string]interface{}
            if err := c.conn.ReadJSON(&msg); err != nil {
                log.Println("Error reading message:", err)
                return
            }
            switch msg["type"] {
            case "init":
                fmt.Println("Document initialized:", msg["content"])
            default:
                fmt.Println("Update received:", msg)
            }
        }
    }()

    // Handle user input
    scanner := bufio.NewScanner(os.Stdin)
    for scanner.Scan() {
        line := scanner.Text()
        parts := strings.SplitN(line, " ", 3)
        if len(parts) < 2 {
            continue
        }
        opType, pos := parts[0], parts[1]
        var char rune
        if opType == "insert" && len(parts) == 3 {
            char = rune(parts[2][0])
        }
        op := Operation{
            Type:     opType,
            Position: atoi(pos),
            Char:     char,
            ClientID: "client1",
            Version:  0,
        }
        c.conn.WriteJSON(op)
    }
}

func atoi(s string) int {
    n, _ := strconv.Atoi(s)
    return n
}

func main() {
    client, err := NewClient("ws://localhost:8080/ws", "doc1")
    if err != nil {
        log.Fatal(err)
    }
    client.Start()
}