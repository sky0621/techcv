package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"

	authinfra "github.com/sky0621/techcv/manager/backend/internal/infrastructure/auth"
	"github.com/sky0621/techcv/manager/backend/internal/infrastructure/clock"
	"github.com/sky0621/techcv/manager/backend/internal/infrastructure/email"
	"github.com/sky0621/techcv/manager/backend/internal/infrastructure/logger"
	"github.com/sky0621/techcv/manager/backend/internal/infrastructure/persistence/memory"
	"github.com/sky0621/techcv/manager/backend/internal/infrastructure/server"
	"github.com/sky0621/techcv/manager/backend/internal/infrastructure/transaction"
	handler "github.com/sky0621/techcv/manager/backend/internal/interface/http/handler"
	httpmiddleware "github.com/sky0621/techcv/manager/backend/internal/interface/http/middleware"
	"github.com/sky0621/techcv/manager/backend/internal/usecase/auth"
	"github.com/sky0621/techcv/manager/backend/internal/usecase/health"
)

func main() {
	log := logger.New()

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	errorHandler := httpmiddleware.NewErrorHandler(log)
	e.HTTPErrorHandler = errorHandler.Handle

	e.Use(echomiddleware.RequestID())
	e.Use(echomiddleware.Recover())
	e.Use(httpmiddleware.Timeout(30 * time.Second))
	e.Use(httpmiddleware.RequestLogger(log))

	clockProvider := clock.NewSystemClock()
	userRepo := memory.NewUserRepository()
	verificationRepo := memory.NewVerificationTokenRepository()
	mailer := email.NewLogMailer(log)
	txManager := transaction.NewNoopManager()
	tokenIssuer := authinfra.NewUUIDTokenIssuer()

	registerConfig := auth.RegisterConfig{
		VerificationURLBase: getEnv("VERIFICATION_URL_BASE", "http://localhost:5173/auth/verify"),
		VerificationTTL:     24 * time.Hour,
	}

	healthUsecase := health.New()
	registerUsecase := auth.NewRegisterUsecase(userRepo, verificationRepo, mailer, clockProvider, registerConfig)
	verifyUsecase := auth.NewVerifyUsecase(userRepo, verificationRepo, txManager, clockProvider, tokenIssuer)
	apiHandler := handler.NewHandler(healthUsecase, registerUsecase, verifyUsecase)

	api := e.Group("/techcv/api/v1")
	apiHandler.Register(api)

	srv := server.New(e, log)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	addr := ":" + getEnv("PORT", "8080")
	log.Info("starting server", "address", addr)

	if err := srv.Start(ctx, addr); err != nil {
		log.Error("server failed", "error", err)
		os.Exit(1)
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
