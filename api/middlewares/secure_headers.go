package middlewares

import (
	"time"

	"github.com/gin-gonic/gin"
)

// SecureHeaders is a middleware that adds common and security-related headers to the HTTP response.
// It sets the 'Connection', 'Content-Type', and 'Date' headers for each request to ensure proper communication settings
// and improve security by controlling caching and connection behaviors.
func SecureHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set the 'Connection' header to 'keep-alive' to maintain persistent connections
		// This allows multiple requests to be sent over the same TCP connection, improving performance.
		c.Header("Connection", "keep-alive")

		// Set the 'Content-Type' header to 'application/json' to specify that the response body is in JSON format
		// This is important for content negotiation and ensuring clients understand the response format.
		c.Header("Content-Type", "application/json")

		// Set the 'Date' header to the current UTC date and time in RFC1123 format
		// This provides clients with information about when the response was generated, useful for caching or debugging.
		c.Header("Date", time.Now().UTC().Format(time.RFC1123))

		// Proceed to the next handler in the chain (or respond to the client if this is the last handler)
		c.Next()
	}
}
