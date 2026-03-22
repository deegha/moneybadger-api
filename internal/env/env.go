package env

import (
	"os"
	"strconv"
)

func GetString(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
}

func GetInt(key string, defaultValue int) int {
	if val := os.Getenv(key); val != "" {
		return atoi(val, defaultValue)
	}
	return defaultValue
}

func atoi(s string, defaultValue int) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return defaultValue
	}
	return i
}
