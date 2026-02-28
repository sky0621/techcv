package server

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/sky0621/techcv/services/manager/backend/internal/infrastructure/config"
)

type Server struct {
	httpServer *http.Server
}

func New(cfg config.Config, handler http.Handler) Server {
	return Server{
		httpServer: &http.Server{
			Addr:    cfg.Addr(),
			Handler: handler,
		},
	}
}

func (s Server) Start() error {
	slog.Info("starting server", "addr", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

func (s Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
