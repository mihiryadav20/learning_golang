package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// LoadConfig loads environment variables from .env file
func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found")
	} else {
		log.Println("Environment variables loaded successfully")
	}
}

// GetEnv retrieves an environment variable or returns a default value if not found
func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// GetJWTSecret returns the JWT secret from environment variables
func GetJWTSecret() []byte {
	secret := GetEnv("JWT_SECRET", "your-256-bit-secret")
	return []byte(secret)
}
