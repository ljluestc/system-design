package api

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

func (s *Server) getHandler(c *gin.Context) {
    key := c.Param("key")
    if !s.scrubber.IsSafe(key) {
        c.JSON(http.StatusForbidden, gin.H{"error": "Malicious request"})
        return
    }
    edge := s.router.Route(key, "mock_location")
    content, ok := edge.Get(key)
    if !ok {
        c.JSON(http.StatusNotFound, gin.H{"error": "Content not found"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"content": content})
}

func (s *Server) putHandler(c *gin.Context) {
    key := c.Param("key")
    var req struct {
        Content string `json:"content"`
    }
    if err := c.BindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
    }
    if !s.scrubber.IsSafe(req.Content) {
        c.JSON(http.StatusForbidden, gin.H{"error": "Malicious content"})
        return
    }
    edge := s.router.Route(key, "mock_location")
    edge.Put(key, req.Content)
    c.JSON(http.StatusOK, gin.H{"status": "success"})
}