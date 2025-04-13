package tests

import (
    "bytes"
    "encoding/json"
    "github.com/calelin/messenger/internal/api"
    "github.com/calelin/messenger/internal/db"
    "github.com/calelin/messenger/pkg/model"
    "github.com/sirupsen/logrus"
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestShortenHandler(t *testing.T) {
    log := logrus.New()
    dbConn := &db.CassandraConn{}
    sequencer := db.NewSequencer(dbConn, log)
    h := api.NewHandler(sequencer, log)

    reqBody, _ := json.Marshal(model.ShortenRequest{OriginalURL: "https://example.com"})
    req, _ := http.NewRequest("POST", "/shorten", bytes.NewReader(reqBody))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("X-API-Key", "test-key")

    rr := httptest.NewRecorder()
    h.ShortenHandler(rr, req)

    if rr.Code != http.StatusOK {
        t.Errorf("Expected status 200, got %v", rr.Code)
    }

    var resp model.ShortenResponse
    if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
        t.Errorf("Failed to decode response: %v", err)
    }
    if !strings.HasPrefix(resp.ShortURL, "http://tiny.url/") {
        t.Errorf("Expected short URL, got %s", resp.ShortURL)
    }
}