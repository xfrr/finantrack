package xos

import "os"

func GetEnvWithDefault(key, defaultValue string) string {
	if value, ok := GetEnv(key); ok {
		return value
	}
	return defaultValue
}

func GetEnv(key string) (string, bool) {
	osValue, ok := os.LookupEnv(key)
	return osValue, ok
}
