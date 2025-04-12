package api

import (
    "encoding/json"
    "net/http"

    "typeahead_system/services"

    "github.com/sirupsen/logrus"
)

// SuggestionAPI manages suggestion requests
type SuggestionAPI struct {
    service *services.SuggestionService
    log     *logrus.Logger
}

// NewSuggestionAPI creates a new SuggestionAPI
func NewSuggestionAPI(service *services.SuggestionService, log *logrus.Logger) *SuggestionAPI {
    return &SuggestionAPI{service: service, log: log}
}

// GetSuggestionsHandler handles GET /suggestions
func (a *SuggestionAPI) GetSuggestionsHandler(w http.ResponseWriter, r *http.Request) {
    prefix := r.URL.Query().Get("prefix")
    if prefix == "" {
        http.Error(w, "Missing prefix parameter", http.StatusBadRequest)
        return
    }
    suggestions, err := a.service.GetSuggestions(prefix, 10) // Limit to 10 suggestions
    if err != nil {
        a.log.Errorf("Failed to get suggestions: %v", err)
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(suggestions)
}