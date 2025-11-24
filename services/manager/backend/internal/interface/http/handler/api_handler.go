// Package handler exposes HTTP handlers that satisfy the OpenAPI contract.
package handler

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/sky0621/techcv/manager/backend/internal/domain"
	httpmiddleware "github.com/sky0621/techcv/manager/backend/internal/interface/http/middleware"
	"github.com/sky0621/techcv/manager/backend/internal/interface/http/response"
	"github.com/sky0621/techcv/manager/backend/internal/usecase/auth"
)

// HealthUsecase defines the behavior required by the handler.
type HealthUsecase interface {
	Check(ctx context.Context) (*domain.HealthStatus, error)
}

// AuthUsecase defines Firebase-backed authentication behaviors.
type AuthUsecase interface {
	Register(ctx context.Context, firebaseUID string) (*auth.AuthResult, error)
	Login(ctx context.Context, firebaseUID string) (*auth.AuthResult, error)
}

// Handler implements the OpenAPI server interface.
type Handler struct {
	health HealthUsecase
	auth   AuthUsecase
}

// NewHandler creates a new API handler instance.
func NewHandler(health HealthUsecase, auth AuthUsecase) *Handler {
	return &Handler{health: health, auth: auth}
}

// Register wires the OpenAPI handlers on the provided Echo group.
func (h *Handler) Register(router *echo.Group, authMiddleware echo.MiddlewareFunc) {
	router.GET("/health", h.GetHealth)

	protected := router.Group("", authMiddleware)
	protected.POST("/auth/firebase/register", h.PostAuthFirebaseRegister)
	protected.POST("/auth/firebase/login", h.PostAuthFirebaseLogin)
}

// GetHealth implements the OpenAPI health endpoint.
func (h *Handler) GetHealth(c echo.Context) error {
	status, err := h.health.Check(c.Request().Context())
	if err != nil {
		return err
	}

	data := map[string]interface{}{
		"status":     status.Status,
		"checked_at": status.CheckedAt,
	}

	return response.Success(c, http.StatusOK, data)
}

// PostAuthFirebaseRegister registers a user after Firebase authentication.
func (h *Handler) PostAuthFirebaseRegister(c echo.Context) error {
	firebaseUID, ok := httpmiddleware.FirebaseUIDFromContext(c)
	if !ok {
		return domain.NewUnauthorized(domain.ErrorCodeInvalidToken, "Firebase UIDが見つかりません", nil)
	}

	result, err := h.auth.Register(c.Request().Context(), firebaseUID)
	if err != nil {
		return err
	}

	return response.Success(c, http.StatusCreated, result)
}

// PostAuthFirebaseLogin logs a user in after Firebase authentication.
func (h *Handler) PostAuthFirebaseLogin(c echo.Context) error {
	firebaseUID, ok := httpmiddleware.FirebaseUIDFromContext(c)
	if !ok {
		return domain.NewUnauthorized(domain.ErrorCodeInvalidToken, "Firebase UIDが見つかりません", nil)
	}

	result, err := h.auth.Login(c.Request().Context(), firebaseUID)
	if err != nil {
		return err
	}

	return response.Success(c, http.StatusOK, result)
}
