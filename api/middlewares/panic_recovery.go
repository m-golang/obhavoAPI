package middlewares

import (
	"fmt"
	"havoAPI/api/helpers"

	"github.com/gin-gonic/gin"
)

// RecoverPanic is a middleware that handles panics in the Gin application.
// If a panic occurs during request processing, it will recover from the panic
// and return a 500 Internal Server Error response to the client.
func RecoverPanic() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Defer function to recover from panic if any occurs during the request lifecycle
		defer func() {
			// Check if a panic occurred (i.e., err is not nil)
			if err := recover(); err != nil {
				// Set the "Connection" header to "close" to indicate the connection should be closed after the response is sent
				c.Header("Connection", "close")

				// Log the error (optional) and send a generic server error response with status 500
				// The error message returned is a generic message, but you could expand this based on the error
				helpers.ServerError(c, fmt.Errorf("this is panic recovery; something went wrong %s", "..."))
			}
		}()

		// Continue processing the request (calls the next handler in the middleware chain)
		c.Next()
	}
}
