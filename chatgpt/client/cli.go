package client

import (
    "flag"
    "fmt"
    "os"
)

type CLI struct {
    client *Client
}

func NewCLI() *CLI {
    return &CLI{client: NewClient("http://localhost:8080")}
}

func (c *CLI) Run() {
    userID := flag.String("user", "user1", "User ID")
    prompt := flag.String("prompt", "", "Prompt")
    flag.Parse()
    if *prompt == "" {
        fmt.Println("Prompt required")
        os.Exit(1)
    }
    resp, err := c.client.Query(*userID, *prompt)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    } else {
        fmt.Printf("Response: %s\n", resp)
    }
}