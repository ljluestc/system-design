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
    cmd := flag.String("cmd", "", "Command: get, put")
    key := flag.String("key", "", "Key")
    content := flag.String("content", "", "Content")
    flag.Parse()

    switch *cmd {
    case "get":
        v, err := c.client.Get(*key)
        if err != nil {
            fmt.Printf("Error: %v\n", err)
        } else {
            fmt.Printf("Content: %s\n", v)
        }
    case "put":
        if err := c.client.Put(*key, *content); err != nil {
            fmt.Printf("Error: %v\n", err)
        } else {
            fmt.Println("Success")
        }
    default:
        fmt.Println("Invalid command")
        os.Exit(1)
    }
}