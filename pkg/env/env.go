package env

import (
	"log"
	"os"
)

// Get returns config from env file
func Get(key string) string {
	env := os.Getenv(key)
	if env == "" {
		log.Fatal("%s is not well-set", key)
	}
	return env
}
