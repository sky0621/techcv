package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sky0621/techcv/services/manager/backend/internal/infrastructure/clock"
	"github.com/sky0621/techcv/services/manager/backend/internal/infrastructure/config"
	"github.com/sky0621/techcv/services/manager/backend/internal/infrastructure/sqlite"
	"github.com/sky0621/techcv/services/manager/backend/internal/interface/http/handler"
	appserver "github.com/sky0621/techcv/services/manager/backend/internal/interface/http/server"
	healthusecase "github.com/sky0621/techcv/services/manager/backend/internal/usecase/health"
)

func main() {
	cfg := config.Load()

	if err := cfg.Validate(); err != nil {
		slog.Error("invalid config", "error", err)
		os.Exit(1)
	}

	if err := sqlite.Prepare(cfg.SQLitePath); err != nil {
		slog.Error("failed to prepare sqlite file", "error", err)
		os.Exit(1)
	}

	healthService := healthusecase.NewService(clock.NewSystemClock())
	healthHandler := handler.NewHealthHandler(healthService)

	mux := http.NewServeMux()
	mux.Handle("GET /health", healthHandler)

	server := appserver.New(cfg, mux)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			slog.Error("server shutdown failed", "error", err)
		}
	}()

	if err := server.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("server stopped with error", "error", err)
		os.Exit(1)
	}
}
