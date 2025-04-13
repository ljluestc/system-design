package auth

import (
    "errors"
    "github.com/google/uuid"
    "messenger/pkg/db"
)

type User struct {
    ID       string
    Email    string
    Password string
    Name     string
}

func (s *Service) CreateUser(email, password, name string) (User, error) {
    // Validate input
    if email == "" || password == "" {
        return User{}, errors.New("email and password required")
    }
    user := User{
        ID:       uuid.New().String(),
        Email:    email,
        Password: password, // In production, hash the password
        Name:     name,
    }
    // Save to database
    if err := s.db.SaveUser(user); err != nil {
        return User{}, err
    }
    return user, nil
}

func (s *Service) GetUser(userID string) (User, error) {
    return s.db.GetUser(userID)
}