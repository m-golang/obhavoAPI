package helpers

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetParametersFromUrl(c *gin.Context) (string, string, error) {
	apiKey := c.Query("key")
	if len(apiKey) == 0 || len(strings.TrimSpace(apiKey)) == 0 {
		return "", "", fmt.Errorf("API key is missing or invalid. Please include a valid API key in your request.")
	}

	query := c.Query("q")
	if len(query) == 0 || len(strings.TrimSpace(query)) == 0 {

		return "", "", fmt.Errorf("Parameter q is missing.")
	}

	return apiKey, query, nil
}
