package health

import (
	"context"
	"testing"
	"time"
)

type fixedClock struct {
	t time.Time
}

func (f fixedClock) Now() time.Time {
	return f.t
}

func TestService_Check(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 2, 28, 9, 30, 0, 0, time.FixedZone("JST", 9*60*60))
	svc := NewService(fixedClock{t: now})

	got := svc.Check(context.Background())

	if got.Service != "manager-backend" {
		t.Fatalf("unexpected service: %s", got.Service)
	}
	if got.Status != "ok" {
		t.Fatalf("unexpected status: %s", got.Status)
	}
	if got.Time != "2026-02-28T00:30:00Z" {
		t.Fatalf("unexpected time: %s", got.Time)
	}
}
