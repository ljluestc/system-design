package main

import (
    "html/template"
    "log"
    "net/http"

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
    r.HandleFunc("/register", RegisterPageHandler).Methods("GET")
    r.HandleFunc("/register", RegisterHandler).Methods("POST")
    r.HandleFunc("/login", LoginPageHandler).Methods("GET")
    r.HandleFunc("/login", LoginHandler).Methods("POST")
    r.HandleFunc("/questions", QuestionHandler).Methods("POST")
    r.HandleFunc("/questions/{id}", QuestionDetailHandler).Methods("GET")
    r.HandleFunc("/search", SearchHandler).Methods("GET")

    // Start the server
    log.Println("Server starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}

// HomeHandler displays the homepage with a list of questions
func HomeHandler(w http.ResponseWriter, r *http.Request) {
    questions, err := GetQuestions()
    if err != nil {
        http.Error(w, "Failed to retrieve questions", http.StatusInternalServerError)
        return
    }
    tmpl := template.Must(template.ParseFiles("templates/home.html"))
    tmpl.Execute(w, questions)
}