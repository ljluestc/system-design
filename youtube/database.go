package main

import (
    "database/sql"
    "log"
    "time"

    _ "github.com/mattn/go-sqlite3"
)

// DB is the global database connection
var DB *sql.DB

// InitDB initializes the SQLite database and creates tables
func InitDB() {
    var err error
    DB, err = sql.Open("sqlite3", "./youtube.db")
    if err != nil {
        log.Fatal(err)
    }

    // Create tables
    _, err = DB.Exec(`
        CREATE TABLE IF NOT EXISTS videos (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            title TEXT,
            description TEXT,
            upload_time DATETIME,
            views INTEGER DEFAULT 0
        );
        CREATE TABLE IF NOT EXISTS comments (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            video_id INTEGER,
            user TEXT,
            comment TEXT,
            timestamp DATETIME
        );
        CREATE TABLE IF NOT EXISTS likes (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            video_id INTEGER,
            user TEXT,
            like_type TEXT  -- "like" or "dislike"
        );
    `)
    if err != nil {
        log.Fatal(err)
    }
}

// Video represents a video entry
type Video struct {
    ID          int
    Title       string
    Description string
    UploadTime  time.Time
    Views       int
}

// InsertVideo adds a new video to the database
func InsertVideo(title, description string) (int64, error) {
    res, err := DB.Exec("INSERT INTO videos (title, description, upload_time) VALUES (?, ?, ?)", title, description, time.Now())
    if err != nil {
        return 0, err
    }
    return res.LastInsertId()
}

// GetVideo retrieves a video by ID
func GetVideo(id int) (*Video, error) {
    row := DB.QueryRow("SELECT id, title, description, upload_time, views FROM videos WHERE id = ?", id)
    v := &Video{}
    err := row.Scan(&v.ID, &v.Title, &v.Description, &v.UploadTime, &v.Views)
    if err != nil {
        return nil, err
    }
    return v, nil
}

// IncrementViews increases the view count for a video
func IncrementViews(id int) error {
    _, err := DB.Exec("UPDATE videos SET views = views + 1 WHERE id = ?", id)
    return err
}

// AddComment adds a comment to a video
func AddComment(videoID int, user, comment string) error {
    _, err := DB.Exec("INSERT INTO comments (video_id, user, comment, timestamp) VALUES (?, ?, ?, ?)", videoID, user, comment, time.Now())
    return err
}

// GetComments retrieves comments for a video
func GetComments(videoID int) ([]map[string]interface{}, error) {
    rows, err := DB.Query("SELECT user, comment, timestamp FROM comments WHERE video_id = ? ORDER BY timestamp DESC", videoID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var comments []map[string]interface{}
    for rows.Next() {
        var user, comment string
        var timestamp time.Time
        err := rows.Scan(&user, &comment, &timestamp)
        if err != nil {
            return nil, err
        }
        comments = append(comments, map[string]interface{}{
            "user":      user,
            "comment":   comment,
            "timestamp": timestamp,
        })
    }
    return comments, nil
}

// AddLike adds a like or dislike to a video
func AddLike(videoID int, user, likeType string) error {
    _, err := DB.Exec("INSERT INTO likes (video_id, user, like_type) VALUES (?, ?, ?)", videoID, user, likeType)
    return err
}

// GetLikes retrieves like and dislike counts for a video
func GetLikes(videoID int) (int, int, error) {
    var likes, dislikes int
    err := DB.QueryRow("SELECT COUNT(*) FROM likes WHERE video_id = ? AND like_type = 'like'", videoID).Scan(&likes)
    if err != nil {
        return 0, 0, err
    }
    err = DB.QueryRow("SELECT COUNT(*) FROM likes WHERE video_id = ? AND like_type = 'dislike'", videoID).Scan(&dislikes)
    if err != nil {
        return 0, 0, err
    }
    return likes, dislikes, nil
}