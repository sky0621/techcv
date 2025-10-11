package server

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// Server wraps an Echo instance and manages graceful shutdown.
type Server struct {
	app    *echo.Echo
	logger *zap.Logger
}

// New constructs a new Server instance.
func New(app *echo.Echo, logger *zap.Logger) *Server {
	return &Server{app: app, logger: logger}
}

// Start runs the HTTP server and listens for shutdown signals via the context.
func (s *Server) Start(ctx context.Context, addr string) error {
	srvErr := make(chan error, 1)

	go func() {
		srvErr <- s.app.Start(addr)
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := s.app.Shutdown(shutdownCtx); err != nil {
			s.logger.Error("graceful shutdown failed", zap.Error(err))
			if errors.Is(err, context.DeadlineExceeded) {
				return s.app.Close()
			}
			return err
		}
		return nil
	case err := <-srvErr:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	}
}
