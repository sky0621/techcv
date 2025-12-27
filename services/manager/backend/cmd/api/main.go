package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"

	"github.com/google/uuid"
	"github.com/sky0621/techcv/manager/backend/internal/infrastructure/clock"
	appconfig "github.com/sky0621/techcv/manager/backend/internal/infrastructure/config"
	"github.com/sky0621/techcv/manager/backend/internal/infrastructure/firebase"
	"github.com/sky0621/techcv/manager/backend/internal/infrastructure/logger"
	"github.com/sky0621/techcv/manager/backend/internal/infrastructure/mysql"
	"github.com/sky0621/techcv/manager/backend/internal/infrastructure/server"
	handler "github.com/sky0621/techcv/manager/backend/internal/interface/http/handler"
	httpmiddleware "github.com/sky0621/techcv/manager/backend/internal/interface/http/middleware"
	"github.com/sky0621/techcv/manager/backend/internal/usecase/auth"
	"github.com/sky0621/techcv/manager/backend/internal/usecase/health"
)

const requestTimeout = 30 * time.Second

func main() {
	cfg, err := appconfig.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	log := logger.New(cfg.App.LogLevel)

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
		if closeErr := db.Close(); closeErr != nil {
			log.Error("failed to close database connection", "error", closeErr)
		}
	}()

	firebaseAuth, err := firebase.NewAuthService(ctx, cfg.Firebase)
	if err != nil {
		log.Error("failed to initialize firebase auth", "error", err)
		os.Exit(1)
	}

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	errorHandler := httpmiddleware.NewErrorHandler(log)
	e.HTTPErrorHandler = errorHandler.Handle

	e.Use(echomiddleware.RequestID())
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.CORSWithConfig(echomiddleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodOptions,
		},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
		},
	}))
	e.Use(httpmiddleware.Timeout(requestTimeout))
	e.Use(httpmiddleware.RequestLogger(log))

	healthRepo := mysql.NewHealthRepository(db)
	healthUsecase := health.New(healthRepo)
	userRepo := mysql.NewUserRepository(db)
	systemClock := clock.NewSystemClock()
	idGenerator := func() (string, error) {
		uid, err := uuid.NewV7()
		if err != nil {
			return "", err
		}
		return uid.String(), nil
	}
	authUsecase := auth.New(firebaseAuth, userRepo, systemClock, idGenerator)
	apiHandler := handler.NewHandler(healthUsecase, authUsecase)

	apiGroup := e.Group("/techcv/api/v1")
	apiHandler.Register(apiGroup, httpmiddleware.FirebaseAuth(firebaseAuth))

	srv := server.New(e, log)

	addr := ":" + cfg.Server.Port
	log.Info("starting server", "address", addr, "env", cfg.App.Environment)

	if err := srv.Start(ctx, addr); err != nil {
		log.Error("server failed", "error", err)
		os.Exit(1)
	}
}
