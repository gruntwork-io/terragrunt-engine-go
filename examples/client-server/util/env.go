// Package util contains helpers shared by the example client and server.
package util

import "os"

func GetEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists { //nolint:forbidigo // terragrunt venv rule; this example reads the process env
		return value
	}

	return defaultValue
}
