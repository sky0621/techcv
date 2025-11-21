// Package handler exposes HTTP handlers that satisfy the OpenAPI contract.
package handler

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/sky0621/techcv/manager/backend/internal/domain"
	openapi "github.com/sky0621/techcv/manager/backend/internal/interface/http/openapi"
	"github.com/sky0621/techcv/manager/backend/internal/interface/http/response"
)

// HealthUsecase defines the behavior required by the handler.
type HealthUsecase interface {
	Check(ctx context.Context) (*domain.HealthStatus, error)
}

// Handler implements the OpenAPI server interface.
type Handler struct {
	health HealthUsecase
}

// NewHandler creates a new API handler instance.
func NewHandler(health HealthUsecase) *Handler {
	return &Handler{health: health}
}

// Register wires the OpenAPI handlers on the provided Echo group.
func (h *Handler) Register(router *echo.Group) {
	openapi.RegisterHandlers(router, h)
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

	meta := map[string]interface{}{
		"requestId": c.Response().Header().Get(echo.HeaderXRequestID),
	}

	return response.Success(c, http.StatusOK, data, meta)
}
