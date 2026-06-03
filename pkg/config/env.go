package config

import (
	"log"
	"os"
)

func required(key string) string {
	value := os.Getenv(key)

	if value == "" {
		log.Fatalf("missing required env: %s", key)
	}

	return value
}

func optional(key, fallback string) string {
	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	return value
}
