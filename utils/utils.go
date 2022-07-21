package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetValue(key string) string {
	var isDevelop bool = os.Getenv("APP_MODE") != "production"

	if isDevelop {
		var err error = godotenv.Load(".env")
		if err != nil {
			log.Fatalf("Error while loading .env file\n")
		}
	}

	return os.Getenv(key)
}
