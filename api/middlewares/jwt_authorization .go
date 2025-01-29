package middlewares

import (
	"fmt"
	"havoAPI/api/config"
	"havoAPI/api/helpers"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// UserAuthorizationJWT checks if the user has a valid JWT token stored in the "u_auth" cookie.
// If the token is missing, invalid, or expired, the request is aborted with an "Unauthorized" response.
// If the token is valid, the userID is extracted from the claims and set in the context for further use by downstream handlers.
func UserAuthorizationJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve the JWT token from the cookie
		tokenStr, err := c.Cookie("u_auth")
		if err != nil {
			helpers.UnauthorizedResponse(c)
			return
		}

		// Parse and validate the JWT token using the signing method and secret key
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			// Ensure the signing method is HMAC (symmetric encryption).
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			// Load the secret key from environment variables to validate the token signature.
			secretKey, err := config.LoadEnvironmentVariable("JWT_SECRET_KEY")
			if err != nil {
				return nil, fmt.Errorf("cannot get secret key while creating and signing JWT: %v", err)
			}

			// Return the secret key for token validation.
			return []byte(secretKey), nil
		})

		// If token parsing or validation fails, return an unauthorized response
		if err != nil || !token.Valid {
			helpers.UnauthorizedResponse(c)
			return
		}

		// Extract claims from the JWT token and check if they are valid
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			helpers.UnauthorizedResponse(c)
			return
		}

		// Check if the token has expired based on the "ttl" claim
		if claims["ttl"].(float64) < float64(time.Now().Unix()) {
			helpers.UnauthorizedResponse(c)
			return
		}

		// Ensure the "userID" claim is valid, otherwise return unauthorized
		userID := claims["userID"].(float64)
		if userID == 0 {
			helpers.UnauthorizedResponse(c)
			return
		}

		// Set the "userID" in the context for further use in downstream handlers.
		c.Set("userID", userID) // Store userID in context.

		// Proceed to the next middleware or handler in the chain.
		c.Next()
	}
}
