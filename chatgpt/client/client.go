package client

import (
    "bytes"
    "encoding/json"
    "net/http"
)

type Client struct {
    baseURL string
}

func NewClient(baseURL string) *Client {
    return &Client{baseURL: baseURL}
}

func (c *Client) Query(userID, prompt string) (string, error) {
    reqBody := map[string]string{"user_id": userID, "prompt": prompt}
    jsonBody, _ := json.Marshal(reqBody)
    resp, err := http.Post(c.baseURL+"/query", "application/json", bytes.NewBuffer(jsonBody))
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()
    var result map[string]string
    json.NewDecoder(resp.Body).Decode(&result)
    return result["response"], nil
}