package main

import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
	"fmt"
)

// URLDatabase manages URL storage
type URLDatabase struct {
    db *sql.DB
}

// NewURLDatabase initializes SQLite database
func NewURLDatabase(path string) (*URLDatabase, error) {
    db, err := sql.Open("sqlite3", path)
    if err != nil {
        return nil, err
    }
    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS urls (
            short_code TEXT PRIMARY KEY,
            long_url TEXT NOT NULL
        )
    `)
    if err != nil {
        return nil, err
    }
    return &URLDatabase{db: db}, nil
}

// Store saves a short code and long URL
func (d *URLDatabase) Store(shortCode, longURL string) error {
    _, err := d.db.Exec("INSERT INTO urls (short_code, long_url) VALUES (?, ?)", shortCode, longURL)
    return err
}

// Retrieve gets the long URL for a short code
func (d *URLDatabase) Retrieve(shortCode string) (string, error) {
    var longURL string
    err := d.db.QueryRow("SELECT long_url FROM urls WHERE short_code = ?", shortCode).Scan(&longURL)
    if err == sql.ErrNoRows {
        return "", fmt.Errorf("not found")
    }
    return longURL, err
}