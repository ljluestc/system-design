package api

import (
    "encoding/json"
    "github.com/calelin/messenger/pkg/db"
    "github.com/calelin/messenger/pkg/model"
    "github.com/gorilla/mux"
    "github.com/sirupsen/logrus"
    "net/http"
)

// DeleteHandler handles DELETE /delete/{shortKey}
// @Summary Delete a short URL
// @Accept json
// @Produce json
// @Param shortKey path string true "Short URL key"
// @Success 200 {string} string "URL Removed"
// @Failure 400 {object} model.ErrorResponse
// @Router /delete/{shortKey} [delete]
func (h *Handler) DeleteHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    shortKey := vars["shortKey"]

    // Validate API key (simplified)
    if err := h.authService.ValidateAPIKey(r.Header.Get("X-API-Key")); err != nil {
        h.log.Errorf("Invalid API key: %v", err)
        respondWithError(w, http.StatusUnauthorized, "Invalid API key")
        return
    }

    // Delete from database
    if err := h.db.DeleteURL(shortKey); err != nil {
        h.log.Errorf("Failed to delete URL %s: %v", shortKey, err)
        respondWithError(w, http.StatusBadRequest, err.Error())
        return
    }

    // Clear cache
    h.cache.Delete(shortKey)
    h.log.Infof("Deleted URL %s", shortKey)
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "URL Removed"})
}