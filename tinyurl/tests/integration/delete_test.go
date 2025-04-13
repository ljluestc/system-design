package integration

import (
    "github.com/calelin/messenger/internal/api"
    "github.com/calelin/messenger/internal/db"
    "github.com/sirupsen/logrus"
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestIntegrationDelete(t *testing.T) {
    log := logrus.New()
    dbConn := &db.CassandraConn{}
    sequencer := db.NewSequencer(dbConn, log)
    h := api.NewHandler(sequencer, log)

    req, _ := http.NewRequest("DELETE", "/delete/abc123", nil)
    req.Header.Set("X-API-Key", "test-key")
    rr := httptest.NewRecorder()
    h.DeleteHandler(rr, req)

    if rr.Code != http.StatusOK {
        t.Errorf("Expected status 200, got %v", rr.Code)
    }
}