package main

import (
    "sync"
    "sync/atomic"
)

// Encoder generates unique short codes
type Encoder struct {
    counter uint64
    mutex   sync.Mutex
}

// NewEncoder initializes the encoder
func NewEncoder() *Encoder {
    return &Encoder{}
}

// Generate creates a new short code using base62 encoding
func (e *Encoder) Generate() string {
    e.mutex.Lock()
    defer e.mutex.Unlock()
    id := atomic.AddUint64(&e.counter, 1)
    return base62Encode(id)
}

// base62Encode converts a number to base62
func base62Encode(num uint64) string {
    const alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
    var encoded string
    for num > 0 {
        encoded = string(alphabet[num%62]) + encoded
        num /= 62
    }
    return encoded
}