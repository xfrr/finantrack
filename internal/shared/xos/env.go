package xos

import (
	"os"
	"strings"
)

func GetEnvWithDefault(key, defaultValue string) string {
	value, ok := GetEnv(key)
	if !ok {
		return defaultValue
	}

	// unquote value if the default value is a string
	if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
		value = strings.Trim(value, "\"")
	}

	return strings.TrimSpace(value)
}

func GetEnv(key string) (string, bool) {
	osValue, ok := os.LookupEnv(key)
	return osValue, ok
}
