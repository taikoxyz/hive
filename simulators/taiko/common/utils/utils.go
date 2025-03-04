package utils

import (
	"os"
	"strings"
)

func GetenvBool(key string) bool {
	val := os.Getenv(key)
	return strings.ToLower(val) == "true"
}
