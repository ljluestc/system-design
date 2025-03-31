package main

import (
    "html/template"
    "log"
    "net/http"
    "strconv"

    "github.com/gorilla/mux"
)

func main() {
    // Initialize the SQLite database
    InitDB()
    defer DB.Close()

    // Set up router
    r := mux.NewRouter()

    // Define routes
    r.HandleFunc("/", HomeHandler).Methods("GET")
    r.HandleFunc("/upload", UploadPageHandler).Methods("GET")
    r.HandleFunc("/upload", UploadHandler).Methods("POST")
    r.HandleFunc("/video/{id}", VideoHandler).Methods("GET")
    r.HandleFunc("/comment", CommentHandler).Methods("POST")
    r.HandleFunc("/like", LikeHandler).Methods("POST")

    // Serve static files (videos and thumbnails)
    r.PathPrefix("/videos/").Handler(http.StripPrefix("/videos/", http.FileServer(http.Dir("./videos"))))
    r.PathPrefix("/thumbnails/").Handler(http.StripPrefix("/thumbnails/", http.FileServer(http.Dir("./thumbnails"))))

    // Start the server
    log.Println("Server starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}

// HomeHandler displays the homepage with a list of videos
func HomeHandler(w http.ResponseWriter, r *http.Request) {
    rows, err := DB.Query("SELECT id, title, description FROM videos ORDER BY upload_time DESC")
    if err != nil {
        http.Error(w, "Failed to retrieve videos", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var videos []Video
    for rows.Next() {
        var v Video
        err := rows.Scan(&v.ID, &v.Title, &v.Description)
        if err != nil {
            http.Error(w, "Failed to scan video", http.StatusInternalServerError)
            return
        }
        videos = append(videos, v)
    }

    tmpl := template.Must(template.ParseFiles("templates/home.html"))
    tmpl.Execute(w, videos)
}

// UploadPageHandler serves the upload form page
func UploadPageHandler(w http.ResponseWriter, r *http.Request) {
    tmpl := template.Must(template.ParseFiles("templates/upload.html"))
    tmpl.Execute(w, nil)
}

// VideoHandler displays a video playback page with comments and likes
func VideoHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid video ID", http.StatusBadRequest)
        return
    }

    video, err := GetVideo(id)
    if err != nil {
        http.Error(w, "Video not found", http.StatusNotFound)
        return
    }

    // Increment view count
    IncrementViews(id)

    comments, err := GetComments(id)
    if err != nil {
        http.Error(w, "Failed to retrieve comments", http.StatusInternalServerError)
        return
    }

    likes, dislikes, err := GetLikes(id)
    if err != nil {
        http.Error(w, "Failed to retrieve likes", http.StatusInternalServerError)
        return
    }

    data := map[string]interface{}{
        "Video":     video,
        "Comments":  comments,
        "Likes":     likes,
        "Dislikes":  dislikes,
    }

    tmpl := template.Must(template.ParseFiles("templates/video.html"))
    tmpl.Execute(w, data)
}

// CommentHandler adds a comment to a video
func CommentHandler(w http.ResponseWriter, r *http.Request) {
    videoID, _ := strconv.Atoi(r.FormValue("video_id"))
    user := r.FormValue("user")
    comment := r.FormValue("comment")
    err := AddComment(videoID, user, comment)
    if err != nil {
        http.Error(w, "Failed to add comment", http.StatusInternalServerError)
        return
    }
    http.Redirect(w, r, "/video/"+strconv.Itoa(videoID), http.StatusSeeOther)
}

// LikeHandler adds a like or dislike to a video
func LikeHandler(w http.ResponseWriter, r *http.Request) {
    videoID, _ := strconv.Atoi(r.FormValue("video_id"))
    user := r.FormValue("user")
    likeType := r.FormValue("like_type")
    err := AddLike(videoID, user, likeType)
    if err != nil {
        http.Error(w, "Failed to add like", http.StatusInternalServerError)
        return
    }
    http.Redirect(w, r, "/video/"+strconv.Itoa(videoID), http.StatusSeeOther)
}