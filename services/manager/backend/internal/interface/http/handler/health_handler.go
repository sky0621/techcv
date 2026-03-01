package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/sky0621/techcv/services/manager/backend/internal/domain/model"
)

type HealthHandler struct {
	healthService interface {
		Check(context.Context) model.HealthStatus
	}
}

func NewHealthHandler(healthService interface {
	Check(context.Context) model.HealthStatus
}) HealthHandler {
	return HealthHandler{healthService: healthService}
}

func (h HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resp := h.healthService.Check(r.Context())

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(resp); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"failed to encode response"}` + "\n"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(buf.Bytes())
}
