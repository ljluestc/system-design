package main

import (
    "html/template"
    "net/http"
    "strconv"

    "github.com/gorilla/mux"
)

// QuestionHandler processes new question submissions
func QuestionHandler(w http.ResponseWriter, r *http.Request) {
    userCookie, err := r.Cookie("user")
    if err != nil || userCookie == nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    title := r.FormValue("title")
    body := r.FormValue("body")

    // Get user ID
    user, err := GetUserByUsername(userCookie.Value)
    if err != nil {
        http.Error(w, "User not found", http.StatusInternalServerError)
        return
    }

    err = InsertQuestion(user.ID, title, body)
    if err != nil {
        http.Error(w, "Failed to insert question", http.StatusInternalServerError)
        return
    }
    http.Redirect(w, r, "/", http.StatusSeeOther)
}

// QuestionDetailHandler displays a question's details
func QuestionDetailHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid question ID", http.StatusBadRequest)
        return
    }

    question, err := GetQuestion(id)
    if err != nil {
        http.Error(w, "Question not found", http.StatusNotFound)
        return
    }

    tmpl := template.Must(template.ParseFiles("templates/question.html"))
    tmpl.Execute(w, question)
}

// SearchHandler handles search queries
func SearchHandler(w http.ResponseWriter, r *http.Request) {
    query := r.URL.Query().Get("query")
    rows, err := DB.Query("SELECT id, title, body FROM questions WHERE title LIKE ? OR body LIKE ?", "%"+query+"%", "%"+query+"%")
    if err != nil {
        http.Error(w, "Search failed", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var questions []Question
    for rows.Next() {
        var q Question
        err := rows.Scan(&q.ID, &q.Title, &q.Body)
        if err != nil {
            http.Error(w, "Failed to scan question", http.StatusInternalServerError)
            return
        }
        questions = append(questions, q)
    }

    tmpl := template.Must(template.ParseFiles("templates/search.html"))
    tmpl.Execute(w, questions)
}