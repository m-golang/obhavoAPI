package config

import (
	"fmt"
	"os"
)
 
// LoadEnvironmentVariable retrieves the value of an environment variable by its key.
// It returns the value of the environment variable as a string if it exists,
// or an error if the variable is not set or is empty.
func LoadEnvironmentVariable(key string) (string, error) {
	// Retrieve the environment variable value using the os.Getenv function.
	value := os.Getenv(key)

	// If the environment variable is empty (not set), return an error.
	if value == "" {
		return "", fmt.Errorf("config: missing environment variable: %s", key)
	}

	// Return the environment variable value if found.
	return value, nil
}
