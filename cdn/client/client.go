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

func (c *Client) Get(key string) (string, error) {
    resp, err := http.Get(c.baseURL + "/content/" + key)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()
    var result map[string]string
    json.NewDecoder(resp.Body).Decode(&result)
    return result["content"], nil
}

func (c *Client) Put(key, content string) error {
    reqBody := map[string]string{"content": content}
    jsonBody, _ := json.Marshal(reqBody)
    _, err := http.Post(c.baseURL+"/content/"+key, "application/json", bytes.NewBuffer(jsonBody))
    return err
}