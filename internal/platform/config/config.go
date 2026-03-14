package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTP     HTTPConfig
	Database DatabaseConfig
	Auth     AuthConfig
}

type HTTPConfig struct {
	Address string
}

type DatabaseConfig struct {
	DSN string
}

type AuthConfig struct {
	APIKey string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := Config{
		HTTP: HTTPConfig{
			Address: ":8080",
		},
	}

	if addr := os.Getenv("HTTP_ADDRESS"); addr != "" {
		cfg.HTTP.Address = addr
	}
	if dsn := os.Getenv("DATABASE_DSN"); dsn != "" {
		cfg.Database.DSN = dsn
	}
	if apiKey := os.Getenv("API_KEY"); apiKey != "" {
		cfg.Auth.APIKey = apiKey
	}

	if cfg.Database.DSN == "" {
		return nil, errors.New("database dsn is required")
	}
	if cfg.Auth.APIKey == "" {
		return nil, errors.New("api key is required")
	}

	return &cfg, nil
}
