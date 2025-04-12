package client

import "net/http"

type Client struct {
    URL string
}

func NewClient(url string) *Client {
    return &Client{URL: url}
}

func (c *Client) Get(key string) (string, error) {
    resp, err := http.Get(c.URL + "/get?key=" + key)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()
    // Placeholder for response parsing
    return "", nil
}