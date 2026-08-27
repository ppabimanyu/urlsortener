package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	GinMode     string
	DatabaseURL string
	JWTSecret   string
	AppBaseURL  string
	FrontendURL string
}

func LoadConfig() *Config {
	// Try to load .env file, ignore error if it doesn't exist in production
	if err := godotenv.Load(); err != nil {
		log.Println("Note: .env file not found, using environment variables")
	}

	return &Config{
		Port:        getEnv("PORT", "8080"),
		GinMode:     getEnv("GIN_MODE", "debug"),
		DatabaseURL: getEnv("DATABASE_URL", "host=localhost user=postgres password=postgrespassword dbname=urlshortener_db port=5432 sslmode=disable"),
		JWTSecret:   getEnv("JWT_SECRET", "default_jwt_secret_key_urlshortener_development"),
		AppBaseURL:  getEnv("APP_BASE_URL", "http://localhost:8080"),
		FrontendURL: getEnv("FRONTEND_URL", "http://localhost:5173"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return fallback
}
