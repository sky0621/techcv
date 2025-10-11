package logger

import (
	"log/slog"
	"os"
	"strings"
)

// New returns a slog logger configured for the current environment.
func New() *slog.Logger {
	env := strings.ToLower(os.Getenv("APP_ENV"))

	var handler slog.Handler

	switch env {
	case "production", "prod":
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	default:
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
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
