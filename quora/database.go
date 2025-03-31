package main

import (
    "database/sql"
    "log"
    "time"

    _ "github.com/mattn/go-sqlite3"
)

// DB is the global database connection
var DB *sql.DB

// InitDB initializes the SQLite database
func InitDB() {
    var err error
    DB, err = sql.Open("sqlite3", "./quora.db")
    if err != nil {
        log.Fatal(err)
    }

    // Create tables
    _, err = DB.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            username TEXT UNIQUE,
            email TEXT UNIQUE,
            password_hash TEXT,
            created_at DATETIME
        );
        CREATE TABLE IF NOT EXISTS questions (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            user_id INTEGER,
            title TEXT,
            body TEXT,
            created_at DATETIME
        );
    `)
    if err != nil {
        log.Fatal(err)
    }
}

// User represents a user entry
type User struct {
    ID           int
    Username     string
    Email        string
    PasswordHash string
    CreatedAt    time.Time
}

// Question represents a question entry
type Question struct {
    ID        int
    UserID    int
    Title     string
    Body      string
    CreatedAt time.Time
}

// InsertUser adds a new user to the database
func InsertUser(username, email, passwordHash string) error {
    _, err := DB.Exec("INSERT INTO users (username, email, password_hash, created_at) VALUES (?, ?, ?, ?)", username, email, passwordHash, time.Now())
    return err
}

// GetUserByUsername retrieves a user by username
func GetUserByUsername(username string) (*User, error) {
    row := DB.QueryRow("SELECT id, username, email, password_hash, created_at FROM users WHERE username = ?", username)
    u := &User{}
    err := row.Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.CreatedAt)
    if err != nil {
        return nil, err
    }
    return u, nil
}

// InsertQuestion adds a new question to the database
func InsertQuestion(userID int, title, body string) error {
    _, err := DB.Exec("INSERT INTO questions (user_id, title, body, created_at) VALUES (?, ?, ?, ?)", userID, title, body, time.Now())
    return err
}

// GetQuestions retrieves all questions
func GetQuestions() ([]Question, error) {
    rows, err := DB.Query("SELECT id, user_id, title, body, created_at FROM questions ORDER BY created_at DESC")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var questions []Question
    for rows.Next() {
        var q Question
        err := rows.Scan(&q.ID, &q.UserID, &q.Title, &q.Body, &q.CreatedAt)
        if err != nil {
            return nil, err
        }
        questions = append(questions, q)
    }
    return questions, nil
}

// GetQuestion retrieves a single question by ID
func GetQuestion(id int) (*Question, error) {
    row := DB.QueryRow("SELECT id, user_id, title, body, created_at FROM questions WHERE id = ?", id)
    q := &Question{}
    err := row.Scan(&q.ID, &q.UserID, &q.Title, &q.Body, &q.CreatedAt)
    if err != nil {
        return nil, err
    }
    return q, nil
}