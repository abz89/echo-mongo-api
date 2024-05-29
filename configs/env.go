package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Get env variable value by key
func GoDotEnvVariable(key string) string {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return os.Getenv(key)
}
