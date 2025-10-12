package handler

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/sky0621/techcv/manager/backend/internal/domain"
	openapi "github.com/sky0621/techcv/manager/backend/internal/interface/http/openapi"
)

// HealthUsecase defines the behaviour required by the handler.
type HealthUsecase interface {
	Check(ctx context.Context) (*domain.HealthStatus, error)
}

// HealthHandler provides HTTP handlers for health endpoints.
type HealthHandler struct {
	usecase HealthUsecase
}

// NewHealthHandler creates a new HealthHandler instance.
func NewHealthHandler(uc HealthUsecase) *HealthHandler {
	return &HealthHandler{usecase: uc}
}

// Register wires the health endpoints on the provided Echo instance.
func (h *HealthHandler) Register(router *echo.Group) {
	openapi.RegisterHandlers(router, h)
}

// GetHealth implements the OpenAPI contract for the health endpoint.
func (h *HealthHandler) GetHealth(c echo.Context) error {
	status, err := h.usecase.Check(c.Request().Context())
	if err != nil {
		return err
	}

	response := openapi.HealthResponseEnvelope{
		Status: "success",
		Data: openapi.HealthStatus{
			Status:    status.Status,
			CheckedAt: status.CheckedAt,
		},
	}

	return c.JSON(http.StatusOK, response)
}
