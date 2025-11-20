// Package config loads application configuration from environment variables.
package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/vrischmann/envconfig"
)

// Config aggregates all runtime configuration for the application.
type Config struct {
	App      AppConfig
	Server   ServerConfig
	Database DatabaseConfig
	Auth     AuthConfig
}

// AppConfig holds application-wide settings.
type AppConfig struct {
	Environment string `envconfig:"APP_ENV,default=development"`
}

// ServerConfig contains HTTP server settings.
type ServerConfig struct {
	Port string `envconfig:"PORT,default=8080"`
}

// DatabaseConfig captures MySQL connection parameters.
type DatabaseConfig struct {
	Host     string `envconfig:"DB_HOST,default=127.0.0.1"`
	Port     string `envconfig:"DB_PORT,default=3306"`
	Name     string `envconfig:"DB_NAME,default=manager"`
	User     string `envconfig:"DB_USER,default=manager"`
	Password string `envconfig:"DB_PASSWORD,default=manager"`
	Params   string `envconfig:"DB_PARAMS,default=parseTime=true&loc=UTC&charset=utf8mb4"`
}

// AuthConfig provides configuration for authentication workflows.
type AuthConfig struct {
	VerificationURLBase string        `envconfig:"VERIFICATION_URL_BASE,default=http://localhost:5173/auth/verify"`
	VerificationTTL     time.Duration `envconfig:"VERIFICATION_TTL,default=24h"`
}

// Load reads .env files (if present) and populates the Config struct via envconfig.
func Load() (*Config, error) {
	if err := loadEnvFiles(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := envconfig.Init(&cfg); err != nil {
		return nil, fmt.Errorf("load configuration: %w", err)
	}

	return &cfg, nil
}

func loadEnvFiles() error {
	files := []string{".env.local", ".env"}
	for _, file := range files {
		if _, err := os.Stat(file); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return fmt.Errorf("stat %s: %w", file, err)
		}

		if err := godotenv.Load(file); err != nil {
			return fmt.Errorf("load %s: %w", file, err)
		}
	}
	return nil
}
