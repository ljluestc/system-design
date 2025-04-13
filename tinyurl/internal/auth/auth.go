package auth

import (
    "github.com/calelin/messenger/pkg/db"
    "github.com/sirupsen/logrus"
)

// AuthService handles authentication
type AuthService struct {
    db  *db.CassandraConn
    log *logrus.Logger
}

// NewAuthService creates a new AuthService
func NewAuthService(db *db.CassandraConn, log *logrus.Logger) *AuthService {
    return &AuthService{db: db, log: log}
}

// ValidateAPIKey checks if an API key is valid
func (s *AuthService) ValidateAPIKey(apiKey string) error {
    // Placeholder: Validate against database
    if apiKey == "" {
        return fmt.Errorf("empty API key")
    }
    s.log.Infof("Validated API key: %s", apiKey)
    return nil
}