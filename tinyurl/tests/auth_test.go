package tests

import (
    "github.com/calelin/messenger/internal/auth"
    "github.com/calelin/messenger/pkg/db"
    "github.com/sirupsen/logrus"
    "testing"
)

func TestValidateAPIKey(t *testing.T) {
    log := logrus.New()
    dbConn := &db.CassandraConn{}
    service := auth.NewAuthService(dbConn, log)
    err := service.ValidateAPIKey("test-key")
    if err != nil {
        t.Errorf("Expected no error, got %v", err)
    }
}