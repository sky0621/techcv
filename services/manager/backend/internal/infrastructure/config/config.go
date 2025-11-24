// Package config loads application configuration from environment variables.
package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/joho/godotenv"
	"github.com/vrischmann/envconfig"
)

var (
	repoRoot     string
	repoRootOnce sync.Once
	envLoaded    []string
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
	if err := loadEnvFiles(); err != nil {
		return nil, err
	}

	if os.Getenv("FIREBASE_PROJECT_ID") == "" {
		return nil, fmt.Errorf("required env FIREBASE_PROJECT_ID not set (loaded env files: %v)", envLoaded)
	}

	var cfg Config
	if err := envconfig.Init(&cfg); err != nil {
		return nil, fmt.Errorf("load configuration: %w (loaded env files: %v)", err, envLoaded)
	}

	return &cfg, nil
}

func loadEnvFiles() error {
	mode := strings.ToLower(strings.TrimSpace(os.Getenv("TECHCV_CONFIG_MODE")))
	if mode == "env-disabled" || mode == "prod" || mode == "production" {
		return nil
	}

	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get working directory: %w", err)
	}

	moduleRoot := findModuleRoot()
	repoRootOnce.Do(func() {
		if _, currentFile, _, ok := runtime.Caller(0); ok {
			// config.go is under internal/infrastructure/config -> move up three levels to reach module root.
			repoRoot = filepath.Clean(filepath.Join(filepath.Dir(currentFile), "../../.."))
		}
	})

	var paths []string
	seen := make(map[string]struct{})
	addPath := func(p string) {
		if p == "" {
			return
		}
		if _, ok := seen[p]; ok {
			return
		}
		paths = append(paths, p)
		seen[p] = struct{}{}
	}

	addPath(filepath.Join(wd, ".env.local"))
	addPath(filepath.Join(wd, ".env"))

	addPath(filepath.Join(moduleRoot, ".env.local"))
	addPath(filepath.Join(moduleRoot, ".env"))

	addPath(filepath.Join(repoRoot, ".env.local"))
	addPath(filepath.Join(repoRoot, ".env"))

	// Also allow .env files placed alongside this config package.
	if _, currentFile, _, ok := runtime.Caller(0); ok {
		configDir := filepath.Dir(currentFile)
		addPath(filepath.Join(configDir, ".env.local"))
		addPath(filepath.Join(configDir, ".env"))
	}

	for _, file := range paths {
		if file == "" {
			continue
		}
		if _, err := os.Stat(file); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return fmt.Errorf("stat %s: %w", file, err)
		}
		if err := godotenv.Overload(file); err != nil {
			return fmt.Errorf("load %s: %w", file, err)
		}
		envLoaded = append(envLoaded, file)
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
