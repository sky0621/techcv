package health

import (
	"context"
	"time"

	"github.com/sky0621/techcv/services/manager/backend/internal/domain/model"
)

type Clock interface {
	Now() time.Time
}

type Service struct {
	clock Clock
}

func NewService(clock Clock) Service {
	return Service{clock: clock}
}

func (s Service) Check(_ context.Context) model.HealthStatus {
	return model.HealthStatus{
		Service: "manager-backend",
		Status:  "ok",
		Time:    s.clock.Now().UTC().Format(time.RFC3339),
	}
}
