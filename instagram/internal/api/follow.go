package api

import (
    "encoding/json"
    "github.com/gorilla/mux"
    "net/http"
)

// FollowHandler handles POST /followUser
func (h *Handler) FollowHandler(w http.ResponseWriter, r *http.Request) {
    var req struct {
        TargetUserID string `json:"target_user_id"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        h.log.Errorf("Invalid request: %v", err)
        respondWithError(w, http.StatusBadRequest, "Invalid request")
        return
    }

    userID := r.Context().Value("userID").(string)
    if err := h.postgresDB.FollowUser(userID, req.TargetUserID); err != nil {
        h.log.Errorf("Failed to follow user: %v", err)
        respondWithError(w, http.StatusInternalServerError, "Failed to follow user")
        return
    }

    h.log.Infof("User %s followed %s", userID, req.TargetUserID)
    w.WriteHeader(http.StatusOK)
}