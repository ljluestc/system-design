package main

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "os"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
)

func TestTinyURLService(t *testing.T) {
    // Initialize service
    service, err := NewTinyURLService()
    assert.NoError(t, err)
    defer os.Remove("tinyurl.db")

    // Start server in a goroutine
    go func() {
        http.HandleFunc("/shorten", service.ShortenHandler)
        http.HandleFunc("/", service.RedirectHandler)
        http.ListenAndServe(":8080", nil)
    }()
    time.Sleep(100 * time.Millisecond) // Wait for server to start

    t.Run("Shorten URL", func(t *testing.T) {
        reqBody := `{"long_url": "https://example.com"}`
        req, _ := http.NewRequest("POST", "http://localhost:8080/shorten", bytes.NewBufferString(reqBody))
        req.Header.Set("Content-Type", "application/json")
        resp, err := http.DefaultClient.Do(req)
        assert.NoError(t, err)
        defer resp.Body.Close()

        assert.Equal(t, http.StatusOK, resp.StatusCode)
        var result map[string]string
        json.NewDecoder(resp.Body).Decode(&result)
        assert.Contains(t, result["short_url"], service.baseURL)
    })

    t.Run("Redirect URL", func(t *testing.T) {
        // Shorten first
        reqBody := `{"long_url": "https://example.com/redirect"}`
        req, _ := http.NewRequest("POST", "http://localhost:8080/shorten", bytes.NewBufferString(reqBody))
        req.Header.Set("Content-Type", "application/json")
        resp, _ := http.DefaultClient.Do(req)
        var result map[string]string
        json.NewDecoder(resp.Body).Decode(&result)
        shortURL := result["short_url"]
        shortCode := shortURL[len(service.baseURL):]

        // Redirect
        req, _ = http.NewRequest("GET", "http://localhost:8080/"+shortCode, nil)
        resp, err = http.DefaultClient.Do(req)
        assert.NoError(t, err)
        defer resp.Body.Close()
        assert.Equal(t, http.StatusFound, resp.StatusCode)
        assert.Equal(t, "https://example.com/redirect", resp.Header.Get("Location"))
    })

    t.Run("Rate Limit", func(t *testing.T) {
        for i := 0; i < 10; i++ {
            req, _ := http.NewRequest("POST", "http://localhost:8080/shorten", bytes.NewBufferString(`{"long_url": "https://example.com"}`))
            req.Header.Set("Content-Type", "application/json")
            resp, _ := http.DefaultClient.Do(req)
            assert.Equal(t, http.StatusOK, resp.StatusCode)
            resp.Body.Close()
        }
        // 11th request should be rate limited
        req, _ := http.NewRequest("POST", "http://localhost:8080/shorten", bytes.NewBufferString(`{"long_url": "https://example.com"}`))
        req.Header.Set("Content-Type", "application/json")
        resp, err := http.DefaultClient.Do(req)
        assert.NoError(t, err)
        assert.Equal(t, http.StatusTooManyRequests, resp.StatusCode)
        resp.Body.Close()
    })

    t.Run("Cache Hit", func(t *testing.T) {
        reqBody := `{"long_url": "https://example.com/cache"}`
        req, _ := http.NewRequest("POST", "http://localhost:8080/shorten", bytes.NewBufferString(reqBody))
        req.Header.Set("Content-Type", "application/json")
        resp1, _ := http.DefaultClient.Do(req)
        var result1 map[string]string
        json.NewDecoder(resp1.Body).Decode(&result1)
        resp1.Body.Close()

        req, _ = http.NewRequest("POST", "http://localhost:8080/shorten", bytes.NewBufferString(reqBody))
        resp2, _ := http.DefaultClient.Do(req)
        var result2 map[string]string
        json.NewDecoder(resp2.Body).Decode(&result2)
        resp2.Body.Close()

        assert.Equal(t, result1["short_url"], result2["short_url"], "Cache should return same short URL")
    })
}