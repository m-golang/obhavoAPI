package helpers

import (
	"fmt"
	"havoAPI/api/config"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)
 
// CreateAndSignJWT generates a JWT token for a given user ID.
// The token includes the user's ID (userID) and an expiration time (ttl).
// The token is signed with a secret key stored in the environment variables.
func CreateAndSignJWT(userID int) (string, error) {
	// Create a new JWT with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":   userID,                                // User ID included in the payload
		"ttl":      time.Now().Add(time.Hour * 24).Unix(), // Token expiration time (1 hour)
	})

	// Load the JWT secret key from environment variables
	secretKey, err := config.LoadEnvironmentVariable("JWT_SECRET_KEY")
	if err != nil {
		return "", fmt.Errorf("cannot get secret key while creating and signing JWT: %v", err)
	}

	// Sign the token with the secret key and return the token string
	return token.SignedString([]byte(secretKey))
}

// SetCookie sets the JWT token as a cookie in the user's browser.
// The cookie is named "u_auth" and will be valid for 1 week (604800 seconds).
// The cookie is marked as HttpOnly for security and will be sent with secure HTTPS connections.
func SetCookie(c *gin.Context, token string) {
	// Set the SameSite attribute for the cookie to Lax, preventing CSRF attacks
	c.SetSameSite(http.SameSiteLaxMode)

	// Set the cookie with the JWT token, with a duration of 1 week (604800 seconds)
	c.SetCookie("u_auth", token, 604800, "", "", false, true)
}

// unauthorizedResponse sends a 401 Unauthorized response with a login prompt message.
// It is used when authentication fails, aborting the request to prevent further processing.
func UnauthorizedResponse(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, gin.H{
		"message": "Please log in to continue",
	})
	c.Abort()
}
