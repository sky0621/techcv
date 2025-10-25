package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"

	authinfra "github.com/sky0621/techcv/manager/backend/internal/infrastructure/auth"
	"github.com/sky0621/techcv/manager/backend/internal/infrastructure/clock"
	"github.com/sky0621/techcv/manager/backend/internal/infrastructure/email"
	"github.com/sky0621/techcv/manager/backend/internal/infrastructure/logger"
	"github.com/sky0621/techcv/manager/backend/internal/infrastructure/mysql"
	"github.com/sky0621/techcv/manager/backend/internal/infrastructure/persistence/memory"
	"github.com/sky0621/techcv/manager/backend/internal/infrastructure/server"
	"github.com/sky0621/techcv/manager/backend/internal/infrastructure/transaction"
	handler "github.com/sky0621/techcv/manager/backend/internal/interface/http/handler"
	httpmiddleware "github.com/sky0621/techcv/manager/backend/internal/interface/http/middleware"
	"github.com/sky0621/techcv/manager/backend/internal/usecase/auth"
	"github.com/sky0621/techcv/manager/backend/internal/usecase/health"
)

const (
	requestTimeout         = 30 * time.Second
	defaultVerificationTTL = 24 * time.Hour
)

type appConfig struct {
	Port                string
	VerificationURLBase string
	DB                  mysql.Config
	Google              googleConfig
	Session             sessionConfig
	Cookie              cookieConfig
}

type googleConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
}

type sessionConfig struct {
	Secret string
	Redis  redisConfig
}

type redisConfig struct {
	Addr     string
	Username string
	Password string
	DB       int
}

type cookieConfig struct {
	Domain string
	Secure bool
}

func main() {
	log := logger.New()

	cfg, err := loadConfig()
	if err != nil {
		log.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	db, err := mysql.NewConnection(ctx, cfg.DB)
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
	clockProvider := clock.NewSystemClock()
	userRepo := memory.NewUserRepository()
	verificationRepo := memory.NewVerificationTokenRepository()
	mailer := email.NewLogMailer(log)
	txManager := transaction.NewNoopManager()
	tokenIssuer := authinfra.NewUUIDTokenIssuer()

	registerConfig := auth.RegisterConfig{
		VerificationURLBase: cfg.VerificationURLBase,
		VerificationTTL:     defaultVerificationTTL,
	}

	registerUsecase := auth.NewRegisterUsecase(userRepo, verificationRepo, mailer, clockProvider, registerConfig)
	verifyUsecase := auth.NewVerifyUsecase(userRepo, verificationRepo, txManager, clockProvider, tokenIssuer)
	apiHandler := handler.NewHandler(healthUsecase, registerUsecase, verifyUsecase)

	apiGroup := e.Group("/techcv/api/v1")
	apiHandler.Register(apiGroup)

	srv := server.New(e, log)

	addr := ":" + cfg.Port
	log.Info("starting server", "address", addr)

	if err := srv.Start(ctx, addr); err != nil {
		log.Error("server failed", "error", err)
		os.Exit(1)
	}
}

func loadConfig() (appConfig, error) {
	cfg := appConfig{
		Port:                getEnv("PORT", "8080"),
		VerificationURLBase: getEnv("VERIFICATION_URL_BASE", "http://localhost:5173/auth/verify"),
		DB: mysql.Config{
			Host:     getEnv("DB_HOST", "127.0.0.1"),
			Port:     getEnv("DB_PORT", "3306"),
			Name:     getEnv("DB_NAME", "manager"),
			User:     getEnv("DB_USER", "manager"),
			Password: getEnv("DB_PASSWORD", "manager"),
			Params:   getEnv("DB_PARAMS", "parseTime=true&loc=UTC&charset=utf8mb4"),
		},
	}

	var err error
	if cfg.Google.ClientID, err = requireEnv("GOOGLE_CLIENT_ID"); err != nil {
		return appConfig{}, err
	}
	if cfg.Google.ClientSecret, err = requireEnv("GOOGLE_CLIENT_SECRET"); err != nil {
		return appConfig{}, err
	}
	if cfg.Google.RedirectURI, err = requireEnv("GOOGLE_REDIRECT_URI"); err != nil {
		return appConfig{}, err
	}

	if cfg.Session.Secret, err = requireEnv("SESSION_SECRET"); err != nil {
		return appConfig{}, err
	}

	if cfg.Cookie.Domain, err = requireEnv("COOKIE_DOMAIN"); err != nil {
		return appConfig{}, err
	}

	if cfg.Cookie.Secure, err = requireBoolEnv("COOKIE_SECURE"); err != nil {
		return appConfig{}, err
	}

	if cfg.Session.Redis.Addr, err = requireEnv("REDIS_ADDR"); err != nil {
		return appConfig{}, err
	}
	cfg.Session.Redis.Username = strings.TrimSpace(os.Getenv("REDIS_USERNAME"))
	cfg.Session.Redis.Password = os.Getenv("REDIS_PASSWORD")

	if cfg.Session.Redis.DB, err = optionalIntEnv("REDIS_DB", 0); err != nil {
		return appConfig{}, err
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func requireEnv(key string) (string, error) {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return "", fmt.Errorf("%s is required", key)
	}
	return value, nil
}

func requireBoolEnv(key string) (bool, error) {
	value, err := requireEnv(key)
	if err != nil {
		return false, err
	}

	parsed, parseErr := strconv.ParseBool(value)
	if parseErr != nil {
		return false, fmt.Errorf("invalid boolean value for %s: %w", key, parseErr)
	}
	return parsed, nil
}

func optionalIntEnv(key string, defaultValue int) (int, error) {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return defaultValue, nil
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("invalid integer value for %s: %w", key, err)
	}
	return parsed, nil
}
