// Package config loads application configuration from environment variables.
package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
	"github.com/vrischmann/envconfig"
)

// Config aggregates all runtime configuration for the application.
type Config struct {
	App      AppConfig
	Server   ServerConfig
	Database DatabaseConfig
	Firebase FirebaseConfig
}

// AppConfig holds application-wide settings.
type AppConfig struct {
	Environment string `envconfig:"TECHCV_APP_ENV,default=development"`
	LogLevel    string `envconfig:"LOG_LEVEL,default=info"`
}

// ServerConfig contains HTTP server settings.
type ServerConfig struct {
	Port string `envconfig:"SERVER_PORT,default=8080"`
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

// FirebaseConfig captures Firebase Authentication settings.
type FirebaseConfig struct {
	ProjectID       string `envconfig:"FIREBASE_PROJECT_ID"`
	CredentialsPath string `envconfig:"FIREBASE_CREDENTIALS_PATH"`
}

// Load reads .env files (if present) and populates the Config struct via envconfig.
func Load() (*Config, error) {
	if err := loadEnvFile(); err != nil {
		return nil, err
	}

	if os.Getenv("FIREBASE_PROJECT_ID") == "" {
		return nil, fmt.Errorf("required env FIREBASE_PROJECT_ID not set")
	}

	var cfg Config
	if err := envconfig.Init(&cfg); err != nil {
		return nil, fmt.Errorf("load configuration: %w", err)
	}

	return &cfg, nil
}

func loadEnvFile() error {
	repoRoot := projectRoot()
	target := filepath.Join(repoRoot, ".env")

	if _, err := os.Stat(target); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("stat %s: %w", target, err)
	}

	if err := godotenv.Overload(target); err != nil {
		return fmt.Errorf("load %s: %w", target, err)
	}
	return nil
}

func findModuleRoot() string {
	search := func(start string) string {
		for dir := start; dir != "" && dir != filepath.Dir(dir); dir = filepath.Dir(dir) {
			if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
				return dir
			}
		}
		return ""
	}

	if _, currentFile, _, ok := runtime.Caller(0); ok {
		if root := search(filepath.Dir(currentFile)); root != "" {
			return root
		}
	}

	if exePath, err := os.Executable(); err == nil {
		if root := search(filepath.Dir(exePath)); root != "" {
			return root
		}
	}

	if wd, err := os.Getwd(); err == nil {
		if root := search(wd); root != "" {
			return root
		}
	}

	return ""
}

func projectRoot() string {
	if root := findModuleRoot(); root != "" {
		return root
	}
	// Fallback to current working directory.
	if wd, err := os.Getwd(); err == nil {
		return wd
	}
	return "."
}
