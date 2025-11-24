// Package logger exposes logging utilities wired for the application.
package logger

import (
	"log/slog"
	"os"
	"strings"
)

// New returns a slog logger configured for the given environment and level.
func New(env, level string) *slog.Logger {
	env = strings.ToLower(env)
	opts := &slog.HandlerOptions{Level: parseLevel(level)}

	var handler slog.Handler
	switch env {
	case "production", "prod":
		handler = slog.NewJSONHandler(os.Stdout, opts)
	default:
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	return slog.New(handler)
}

// WithRequestID enriches the logger with a request identifier when present.
func WithRequestID(base *slog.Logger, requestID string) *slog.Logger {
	if requestID == "" {
		return base
	}
	return base.With("request_id", requestID)
}

func parseLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
