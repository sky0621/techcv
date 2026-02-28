package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Env        string
	Port       string
	SQLitePath string
}

func Load() Config {
	return Config{
		Env:        envOrDefault("APP_ENV", "local"),
		Port:       envOrDefault("APP_PORT", "8080"),
		SQLitePath: envOrDefault("SQLITE_PATH", "./tmp/manager.db"),
	}
}

func (c Config) Addr() string {
	return ":" + c.Port
}

func (c Config) Validate() error {
	port, err := strconv.Atoi(c.Port)
	if err != nil {
		return fmt.Errorf("APP_PORT must be numeric: %w", err)
	}
	if port < 1 || port > 65535 {
		return fmt.Errorf("APP_PORT must be in range 1-65535: %d", port)
	}
	return nil
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
