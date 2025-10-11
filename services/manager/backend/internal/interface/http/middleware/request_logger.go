package middleware

import (
	"log/slog"
	"time"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"

	"github.com/sky0621/techcv/manager/backend/internal/infrastructure/logger"
)

// RequestLogger logs basic request/response information using the shared logger.
func RequestLogger(log *slog.Logger) echo.MiddlewareFunc {
	cfg := echomiddleware.RequestLoggerConfig{
		LogURI:     true,
		LogStatus:  true,
		LogLatency: true,
		LogMethod:  true,
		LogError:   true,
		LogValuesFunc: func(c echo.Context, v echomiddleware.RequestLoggerValues) error {
			requestID := c.Response().Header().Get(echo.HeaderXRequestID)
			entry := logger.WithRequestID(log, requestID)

			if v.Error == nil {
				entry.Info("request handled",
					slog.String("method", v.Method),
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
					slog.Duration("latency", v.Latency),
				)
				return nil
			}

			entry.Error("request failed",
				slog.String("method", v.Method),
				slog.String("uri", v.URI),
				slog.Int("status", v.Status),
				slog.Duration("latency", v.Latency),
				slog.Any("error", v.Error),
			)
			return nil
		},
	}

	return echomiddleware.RequestLoggerWithConfig(cfg)
}

// Timeout applies a simple request timeout.
func Timeout(duration time.Duration) echo.MiddlewareFunc {
	return echomiddleware.TimeoutWithConfig(echomiddleware.TimeoutConfig{
		Timeout: duration,
	})
}
