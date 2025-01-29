package config

import (
	"fmt"
	"os"
)

func LoadEnvironmentVariable(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", fmt.Errorf("config: missing environment variable: %s", key)
	}

	return value, nil
}
