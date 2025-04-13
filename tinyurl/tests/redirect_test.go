package tests

import (
    "github.com/calelin/messenger/internal/api"
    "github.com/calelin/messenger/internal/db"
    "github.com/sirupsen/logrus"
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestRedirectHandler(t *testing.T) {
    log := logrus.New()
    dbConn := &db.CassandraConn{}
    sequencer := db.NewSequencer(dbConn, log)
    h := api.NewHandler(sequencer, log)

    req, _ := http.NewRequest("GET", "/abc123", nil)
    rr := httptest.NewRecorder()
    h.RedirectHandler(rr, req)

    if rr.Code != http.StatusNotFound {
        t.Errorf("Expected status 404, got %v", rr.Code)
    }
}