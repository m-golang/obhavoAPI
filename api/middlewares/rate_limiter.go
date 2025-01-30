package middlewares

import (
	"havoAPI/api/helpers"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func RateLimiter() gin.HandlerFunc {
	limiter := rate.NewLimiter(10, 30)

	return func(c *gin.Context) {
		if !limiter.Allow() {
			helpers.RateLimitExceededResponse(c)
			return
		}

		c.Next()
	}
}
