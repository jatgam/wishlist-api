package utils

import (
	"os"
	"strconv"
)

// GetEnv will return the value of an environment variable if it exists
// otherwise it will return the supplied defaultValue
func GetEnv(key string, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = defaultValue
	}
	return value
}

// GetEnvAsBool will return an environment var if it exists, parsed as a bool
// If parsing fails, false is returned. Otherwise the default value is returned
func GetEnvAsBool(key string, defaultValue string) bool {
	value, err := strconv.ParseBool(GetEnv(key, defaultValue))
	if err != nil {
		return false
	}
	return value
}
