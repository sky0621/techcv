package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/sky0621/techcv/services/manager/backend/internal/domain"
)

type HealthHandler struct {
	healthService interface {
		Check(context.Context) domain.HealthStatus
	}
}

func NewHealthHandler(healthService interface {
	Check(context.Context) domain.HealthStatus
}) HealthHandler {
	return HealthHandler{healthService: healthService}
}

func (h HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := h.healthService.Check(r.Context())
	_ = json.NewEncoder(w).Encode(resp)
}
