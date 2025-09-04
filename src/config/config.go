package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName      string
	Host         string
	Port         string
	DevMode      bool
	AdminName    string
	AdminEmail   string
	AdminPassword string
	BotToken     string
	BotUsername  string
	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
	SMTPFrom     string
	DBPath       string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		if err := godotenv.Load("build/.env"); err != nil {
			return nil, fmt.Errorf("failed to load .env file: %w", err)
		}
	}

	devMode := getEnv("DEV_MODE", "false") == "true"
	var host string
	if devMode {
		// In dev mode, use localhost with the configured port
		port := getEnv("PORT", "8080")
		host = "http://localhost:" + port
	} else {
		// In production mode, use the HOST env var
		host = getEnv("HOST", "http://localhost:8080")
	}

	cfg := &Config{
		AppName:       getEnv("APP_NAME", "Disciplo"),
		Host:          host,
		Port:          getEnv("PORT", "8080"),
		DevMode:       devMode,
		AdminName:     getEnv("ADMIN_NAME", "Admin"),
		AdminEmail:    getEnvRequired("ADMIN_EMAIL"),
		AdminPassword: getEnvRequired("ADMIN_PASSWORD"),
		BotToken:      getEnvRequired("BOT_TOKEN"),
		BotUsername:   getEnv("BOT_USERNAME", ""),
		SMTPHost:      getEnv("SMTP_HOST", ""),
		SMTPPort:      getEnv("SMTP_PORT", "587"),
		SMTPUsername:  getEnv("SMTP_USER", getEnv("SMTP_USERNAME", "")),
		SMTPPassword:  getEnv("SMTP_PASS", getEnv("SMTP_PASSWORD", "")),
		SMTPFrom:      getEnv("SMTP_FROM", ""),
		DBPath:        getEnv("DB_PATH", "pb_data"),
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvRequired(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("Required environment variable %s is not set", key))
	}
	return value
}