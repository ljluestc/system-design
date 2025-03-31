package main

import (
    "crypto/sha256"
    "encoding/hex"
    "html/template"
    "net/http"
)

// RegisterPageHandler serves the registration page
func RegisterPageHandler(w http.ResponseWriter, r *http.Request) {
    tmpl := template.Must(template.ParseFiles("templates/register.html"))
    tmpl.Execute(w, nil)
}

// RegisterHandler processes user registration
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
    username := r.FormValue("username")
    email := r.FormValue("email")
    password := r.FormValue("password")

    // Hash password
    hash := sha256.New()
    hash.Write([]byte(password))
    passwordHash := hex.EncodeToString(hash.Sum(nil))

    err := InsertUser(username, email, passwordHash)
    if err != nil {
        http.Error(w, "Failed to register user", http.StatusInternalServerError)
        return
    }
    http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// LoginPageHandler serves the login page
func LoginPageHandler(w http.ResponseWriter, r *http.Request) {
    tmpl := template.Must(template.ParseFiles("templates/login.html"))
    tmpl.Execute(w, nil)
}

// LoginHandler processes user login
func LoginHandler(w http.ResponseWriter, r *http.Request) {
    username := r.FormValue("username")
    password := r.FormValue("password")

    user, err := GetUserByUsername(username)
    if err != nil || user == nil {
        http.Error(w, "Invalid credentials", http.StatusUnauthorized)
        return
    }

    // Verify password
    hash := sha256.New()
    hash.Write([]byte(password))
    if hex.EncodeToString(hash.Sum(nil)) != user.PasswordHash {
        http.Error(w, "Invalid credentials", http.StatusUnauthorized)
        return
    }

    // Set session cookie (simplified)
    http.SetCookie(w, &http.Cookie{Name: "user", Value: username})
    http.Redirect(w, r, "/", http.StatusSeeOther)
}