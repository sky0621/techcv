package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"

	"github.com/sky0621/techcv/manager/backend/internal/infrastructure/logger"
)

// RequestLogger logs basic request/response information using the shared logger.
func RequestLogger(log *zap.Logger) echo.MiddlewareFunc {
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
					zap.String("method", v.Method),
					zap.String("uri", v.URI),
					zap.Int("status", v.Status),
					zap.Duration("latency", v.Latency),
				)
				return nil
			}

			entry.Error("request failed",
				zap.String("method", v.Method),
				zap.String("uri", v.URI),
				zap.Int("status", v.Status),
				zap.Duration("latency", v.Latency),
				zap.Error(v.Error),
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
