package handler

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/vibe-kanban/backend/internal/domain"
	"github.com/vibe-kanban/backend/internal/interface/http/response"
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
	router.GET("/health", h.check)
}

func (h *HealthHandler) check(c echo.Context) error {
	status, err := h.usecase.Check(c.Request().Context())
	if err != nil {
		return err
	}
	return response.Success(c, http.StatusOK, status)
}
