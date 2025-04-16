package api

import (
    "cdn/routing"
    "cdn/scrubber"
    "github.com/gin-gonic/gin"
)

type Server struct {
    router   *routing.Router
    scrubber *scrubber.Service
}

func NewServer(router *routing.Router, scrubber *scrubber.Service) *Server {
    return &Server{router: router, scrubber: scrubber}
}

func (s *Server) Start() error {
    r := gin.Default()
    r.GET("/content/:key", s.getHandler)
    r.PUT("/content/:key", s.putHandler)
    return r.Run(":8080")
}