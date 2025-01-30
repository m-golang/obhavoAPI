package middlewares

import (
	"time"

	"github.com/gin-gonic/gin"
)

// SecureHeaders is a middleware that adds secure and common headers to the HTTP response.
// It sets the 'Connection', 'Content-Type', and 'Date' headers for each request.
func SecureHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Connection", "keep-alive")
		c.Header("Content-Type", "application/json")
		c.Header("Date", time.Now().UTC().Format(time.RFC1123))

		c.Next()
	}
}
