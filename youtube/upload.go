package main

import (
    "fmt"
    "io"
    "net/http"
    "os"
    "path/filepath"
    "strconv"
)

// UploadHandler processes video uploads
func UploadHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    // Parse multipart form (max 32MB)
    err := r.ParseMultipartForm(32 << 20)
    if err != nil {
        http.Error(w, "Unable to parse form", http.StatusBadRequest)
        return
    }

    // Get form values
    title := r.FormValue("title")
    description := r.FormValue("description")

    // Get video file
    file, handler, err := r.FormFile("video")
    if err != nil {
        http.Error(w, "Unable to get video file", http.StatusBadRequest)
        return
    }
    defer file.Close()

    // Insert video metadata into database
    id, err := InsertVideo(title, description)
    if err != nil {
        http.Error(w, "Failed to insert video", http.StatusInternalServerError)
        return
    }

    // Save video file
    videoPath := filepath.Join("videos", fmt.Sprintf("%d%s", id, filepath.Ext(handler.Filename)))
    out, err := os.Create(videoPath)
    if err != nil {
        http.Error(w, "Unable to save video", http.StatusInternalServerError)
        return
    }
    defer out.Close()
    _, err = io.Copy(out, file)
    if err != nil {
        http.Error(w, "Unable to save video", http.StatusInternalServerError)
        return
    }

    // Generate placeholder thumbnail
    thumbPath := filepath.Join("thumbnails", fmt.Sprintf("%d.jpg", id))
    err = os.WriteFile(thumbPath, []byte("placeholder"), 0644)
    if err != nil {
        http.Error(w, "Unable to save thumbnail", http.StatusInternalServerError)
        return
    }

    // Redirect to homepage
    http.Redirect(w, r, "/", http.StatusSeeOther)
}