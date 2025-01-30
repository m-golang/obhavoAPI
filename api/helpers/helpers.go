package helpers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ServerError logs unexpected server errors and returns a generic internal server error response.
// It ensures sensitive information about the error is not exposed to the client.
func ServerError(c *gin.Context, err error) {
	// Log the error on the server for further inspection
	log.Println(err)
	// Send a generic error response to the client
	c.JSON(http.StatusInternalServerError, gin.H{
		"error": "An unexpected server error occurred. Please try again later.",
	})
}

// ClientError is used to respond to client errors with a specific HTTP status code and message.
// This can be used for known client-side errors (e.g., 400 for bad requests, 404 for not found).
func ClientError(c *gin.Context, code int, message string) {
	// Return the client error response with the provided status code and message
	c.JSON(code, gin.H{
		"error": message,
	})
}

