package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"

	appconfig "github.com/sky0621/techcv/manager/backend/internal/infrastructure/config"
	"github.com/sky0621/techcv/manager/backend/internal/infrastructure/logger"
	"github.com/sky0621/techcv/manager/backend/internal/infrastructure/mysql"
	"github.com/sky0621/techcv/manager/backend/internal/infrastructure/server"
	handler "github.com/sky0621/techcv/manager/backend/internal/interface/http/handler"
	httpmiddleware "github.com/sky0621/techcv/manager/backend/internal/interface/http/middleware"
	"github.com/sky0621/techcv/manager/backend/internal/usecase/health"
)

const requestTimeout = 30 * time.Second

func main() {
	cfg, err := appconfig.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	log := logger.New()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	dbCfg := mysql.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		Name:     cfg.Database.Name,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		Params:   cfg.Database.Params,
	}

	db, err := mysql.NewConnection(ctx, dbCfg)
	if err != nil {
		log.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Error("failed to close database connection", "error", err)
		}
	}()

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	errorHandler := httpmiddleware.NewErrorHandler(log)
	e.HTTPErrorHandler = errorHandler.Handle

	e.Use(echomiddleware.RequestID())
	e.Use(echomiddleware.Recover())
	e.Use(httpmiddleware.Timeout(requestTimeout))
	e.Use(httpmiddleware.RequestLogger(log))

	healthRepo := mysql.NewHealthRepository(db)
	healthUsecase := health.New(healthRepo)
	apiHandler := handler.NewHandler(healthUsecase)

	apiGroup := e.Group("/techcv/api/v1")
	apiHandler.Register(apiGroup)

	srv := server.New(e, log)

	addr := ":" + cfg.Server.Port
	log.Info("starting server", "address", addr, "env", cfg.App.Environment)

	if err := srv.Start(ctx, addr); err != nil {
		log.Error("server failed", "error", err)
		os.Exit(1)
	}
}
