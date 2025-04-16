package api

import (
    "chatgpt/moderation"
    "net/http"

    "github.com/gin-gonic/gin"
)

func (s *Server) queryHandler(c *gin.Context) {
    var req struct {
        UserID string `json:"user_id"`
        Prompt string `json:"prompt"`
    }
    if err := c.BindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
    }
    filter := moderation.NewFilter()
    if !filter.IsSafe(req.Prompt) {
        c.JSON(http.StatusForbidden, gin.H{"error": "Unsafe prompt"})
        return
    }
    response, err := s.router.ProcessQuery(req.UserID, req.Prompt)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"response": response})
}