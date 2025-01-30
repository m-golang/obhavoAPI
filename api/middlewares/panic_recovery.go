package middlewares

import (
	"fmt"
	"havoAPI/api/helpers"

	"github.com/gin-gonic/gin"
)

// RecoverPanic is a middleware that handles panics in the Gin application.
// If a panic occurs during request processing, it will recover and return a 500 Internal Server Error response.
func RecoverPanic() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Defer function to recover from panic if any occurs during the request lifecycle
		defer func() {
			if err := recover(); err != nil {
				c.Header("Connection", "close")
				helpers.ServerError(c, fmt.Errorf("this is panic recovery; something went wrong %s","..."))
			}
		}()
		c.Next()
	}
}
