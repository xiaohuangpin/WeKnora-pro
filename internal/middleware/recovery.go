package middleware

import (
	"fmt"
	"log"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

// Recovery is a middleware that recovers from panics
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Get request ID
				requestID, _ := c.Get("RequestID")

				// Print stacktrace
				stacktrace := debug.Stack()
				// Log error
				log.Printf("[PANIC] %s | %v | %s", requestID, err, stacktrace)

				// 返回500错误
				c.AbortWithStatusJSON(500, gin.H{
					"error":   "Internal Server Error",
					"message": fmt.Sprintf("%v", err),
				})
			}
		}()

		c.Next()
	}
}
