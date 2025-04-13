package api

import (
    "github.com/calelin/messenger/internal/cache"
    "github.com/calelin/messenger/pkg/db"
    "github.com/sirupsen/logrus"
    "net/http"
)

// RedirectHandler handles GET /{shortKey}
// @Summary Redirect to original URL
// @Produce json
// @Param shortKey path string true "Short URL key"
// @Success 302 {string} string "Redirect to original URL"
// @Failure 404 {object} model.ErrorResponse
// @Router /{shortKey} [get]
func (h *Handler) RedirectHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    shortKey := vars["shortKey"]

    // Check cache
    originalURL, err := h.cache.Get(shortKey)
    if err == nil {
        http.Redirect(w, r, originalURL, http.StatusFound)
        return
    }

    // Check database
    originalURL, err = h.db.GetOriginalURL(shortKey)
    if err != nil {
        h.log.Errorf("Failed to find URL for %s: %v", shortKey, err)
        respondWithError(w, http.StatusNotFound, "URL not found")
        return
    }

    // Cache the result
    h.cache.Set(shortKey, originalURL)
    http.Redirect(w, r, originalURL, http.StatusFound)
}