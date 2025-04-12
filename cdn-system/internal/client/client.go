package client

import (
    "fmt"
    "log"
    "net/http"
)

// Client interacts with the CDN
type Client struct {
    URL string
}

// NewClient creates a client instance
func NewClient(url string) *Client {
    return &Client{URL: url}
}

// GetContent fetches content
func (c *Client) GetContent(key string) (string, error) {
    resp, err := http.Get(c.URL + "/content?key=" + key)
    if err != nil {
        return "", fmt.Errorf("request failed: %v", err)
    }
    defer resp.Body.Close()
    log.Printf("Fetched content for key %s", key)
    return "placeholder", nil // Simplified
}

// Expand with more client methods
// ...