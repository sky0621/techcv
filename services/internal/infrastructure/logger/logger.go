package logger

import (
	"os"
	"strings"

	"go.uber.org/zap"
)

// New returns a zap logger configured for the current environment.
func New() *zap.Logger {
	env := strings.ToLower(os.Getenv("APP_ENV"))

	switch env {
	case "production", "prod":
		logger, err := zap.NewProduction()
		if err != nil {
			panic(err)
		}
		return logger
	default:
		logger, err := zap.NewDevelopment()
		if err != nil {
			panic(err)
		}
		return logger
	}
}

// WithRequestID enriches the logger with a request identifier when present.
func WithRequestID(base *zap.Logger, requestID string) *zap.Logger {
	if requestID == "" {
		return base
	}
	return base.With(zap.String("request_id", requestID))
}
