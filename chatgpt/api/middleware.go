package api

import "github.com/gin-gonic/gin"

func RateLimitMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next() // Placeholder
    }
}

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next() // Placeholder
    }
}