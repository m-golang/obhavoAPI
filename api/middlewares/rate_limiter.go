package middlewares

import (
	"havoAPI/api/helpers"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimiter is a middleware that limits the number of requests that can be made in a given time window.
// It uses a token bucket algorithm to allow a certain number of requests (10 requests) every 30 seconds.
// If the rate limit is exceeded, it responds with a 429 Too Many Requests status.
func RateLimiter() gin.HandlerFunc {
	// Create a new rate limiter with a maximum of 10 requests per 30 seconds.
	limiter := rate.NewLimiter(10, 30)

	return func(c *gin.Context) {
		// Check if the current request is allowed based on the rate limit
		if !limiter.Allow() {
			// If the rate limit is exceeded, return a rate limit exceeded response
			helpers.RateLimitExceededResponse(c)
			return
		}

		// If the request is allowed, proceed to the next middleware or handler in the chain
		c.Next()
	}
}
