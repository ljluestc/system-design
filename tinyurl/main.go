// Copyright 2025 Your Name.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
    "encoding/json"
    "log"
    "net/http"
    "strings"
    "sync"
    "time"
)

// TinyURLService encapsulates the service components
type TinyURLService struct {
    db         *URLDatabase
    encoder    *Encoder
    cache      *LRUCache
    rateLimiter *RateLimiter
    baseURL    string
    mu         sync.Mutex
}

// NewTinyURLService initializes the TinyURL service
func NewTinyURLService() (*TinyURLService, error) {
    db, err := NewURLDatabase("tinyurl.db")
    if err != nil {
        return nil, err
    }
    return &TinyURLService{
        db:         db,
        encoder:    NewEncoder(),
        cache:      NewLRUCache(1000), // 1000 items max
        rateLimiter: NewRateLimiter(10, 60*time.Second), // 10 reqs/sec
        baseURL:    "http://tinyurl.local/",
    }, nil
}

// ShortenHandler handles POST /shorten requests
func (s *TinyURLService) ShortenHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    if !s.rateLimiter.Allow() {
        http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
        return
    }

    var req struct {
        LongURL string `json:"long_url"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    s.mu.Lock()
    defer s.mu.Unlock()

    // Check cache first
    if shortURL, ok := s.cache.Get(req.LongURL); ok {
        jsonResponse(w, map[string]string{"short_url": shortURL})
        return
    }

    // Generate short code and store
    shortCode := s.encoder.Generate()
    shortURL := s.baseURL + shortCode
    if err := s.db.Store(shortCode, req.LongURL); err != nil {
        http.Error(w, "Failed to store URL", http.StatusInternalServerError)
        return
    }

    // Update cache
    s.cache.Set(req.LongURL, shortURL)

    jsonResponse(w, map[string]string{"short_url": shortURL})
}

// RedirectHandler handles GET /{shortCode} requests
func (s *TinyURLService) RedirectHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    if !s.rateLimiter.Allow() {
        http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
        return
    }

    shortCode := strings.TrimPrefix(r.URL.Path, "/")
    longURL, err := s.db.Retrieve(shortCode)
    if err != nil {
        http.Error(w, "URL not found", http.StatusNotFound)
        return
    }

    http.Redirect(w, r, longURL, http.StatusFound)
}

// jsonResponse sends a JSON response
func jsonResponse(w http.ResponseWriter, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(data); err != nil {
        http.Error(w, "Failed to encode response", http.StatusInternalServerError)
    }
}

func main() {
    service, err := NewTinyURLService()
    if err != nil {
        log.Fatalf("Failed to initialize service: %v", err)
    }

    http.HandleFunc("/shorten", service.ShortenHandler)
    http.HandleFunc("/", service.RedirectHandler)

    log.Println("Starting TinyURL server on :8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatalf("Server failed: %v", err)
    }
}