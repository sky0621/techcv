package config

import "os"

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

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
