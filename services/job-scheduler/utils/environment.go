package utils

import (
	"fmt"
	"os"
)

// Loads an Environment variable as required, returning an error if not found
func LoadEnvRequired(name string) (string, error) {
	value := os.Getenv(name)
	if value == "" {
		return "", fmt.Errorf("environment variable %s is required", name)
	}
	return value, nil
}

// Loads an Environment variable or returns a default value if not found
func LoadEnvOrDefault(name string, defaultValue string) string {
	value := os.Getenv(name)
	if value == "" {
		return defaultValue
	}
	return value
}
