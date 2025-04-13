package model

import "time"

// ShortenRequest represents a shorten request
type ShortenRequest struct {
    OriginalURL  string     `json:"original_url"`
    CustomAlias  string     `json:"custom_alias,omitempty"`
    ExpiryDate   *time.Time `json:"expiry_date,omitempty"`
}

// ShortenResponse represents a shorten response
type ShortenResponse struct {
    ShortURL string `json:"short_url"`
}