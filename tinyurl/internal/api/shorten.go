package api

import (
    "encoding/json"
    "github.com/calelin/tinyurl/internal/db"
    "github.com/calelin/tinyurl/pkg/model"
    "github.com/sirupsen/logrus"
    "net/http"
)

// Handler manages API endpoints
type Handler struct {
    sequencer *db.Sequencer
    log       *logrus.Logger
}

// NewHandler creates a new Handler
func NewHandler(sequencer *db.Sequencer, log *logrus.Logger) *Handler {
    return &Handler{sequencer: sequencer, log: log}
}

// ShortenHandler handles POST /shorten
// @Summary Shorten a URL
// @Accept json
// @Produce json
// @Param request body model.ShortenRequest true "Shorten request"
// @Success 200 {object} model.ShortenResponse
// @Failure 400 {object} model.ErrorResponse
// @Router /shorten [post]
func (h *Handler) ShortenHandler(w http.ResponseWriter, r *http.Request) {
    var req model.ShortenRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        h.log.Errorf("Invalid request: %v", err)
        respondWithError(w, http.StatusBadRequest, "Invalid request")
        return
    }

    // Generate short URL
    shortURL, err := h.sequencer.GenerateShortURL(req.OriginalURL, req.CustomAlias, req.ExpiryDate)
    if err != nil {
        h.log.Errorf("Failed to shorten URL: %v", err)
        respondWithError(w, http.StatusBadRequest, err.Error())
        return
    }

    resp := model.ShortenResponse{ShortURL: shortURL}
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp)
}

// respondWithError sends an error response
func respondWithError(w http.ResponseWriter, code int, message string) {
    resp := model.ErrorResponse{Error: message}
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    json.NewEncoder(w).Encode(resp)
}