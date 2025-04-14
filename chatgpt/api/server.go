package api

import (
    "chatgpt/inference"
    "github.com/gin-gonic/gin"
)

type Server struct {
    router *inference.Router
}

func NewServer(router *inference.Router) *Server {
    return &Server{router: router}
}

func (s *Server) Start() error {
    r := gin.Default()
    r.POST("/query", s.queryHandler)
    return r.Run(":8080")
}