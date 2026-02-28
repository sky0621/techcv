package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sky0621/techcv/services/manager/backend/internal/domain"
	profileregister "github.com/sky0621/techcv/services/manager/backend/internal/usecase/profile/register"
)

type stubProfileRegistrar struct {
	registerFn func(ctx context.Context, name, nickname string) (domain.Profile, error)
}

func (s stubProfileRegistrar) Register(ctx context.Context, name, nickname string) (domain.Profile, error) {
	return s.registerFn(ctx, name, nickname)
}

func TestProfileHandler_ServeHTTP_Created(t *testing.T) {
	t.Parallel()

	h := NewProfileHandler(stubProfileRegistrar{
		registerFn: func(_ context.Context, name, nickname string) (domain.Profile, error) {
			return domain.Profile{
				ID:       "profile_1",
				Name:     name,
				Nickname: nickname,
			}, nil
		},
	})

	req := httptest.NewRequest(http.MethodPost, "/profiles", bytes.NewBufferString(`{"name":"Alice","nickname":"ali"}`))
	res := httptest.NewRecorder()

	h.ServeHTTP(res, req)

	if res.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusCreated)
	}
	var body domain.Profile
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Name != "Alice" || body.Nickname != "ali" {
		t.Fatalf("unexpected response: %+v", body)
	}
}

func TestProfileHandler_ServeHTTP_BadRequest(t *testing.T) {
	t.Parallel()

	h := NewProfileHandler(stubProfileRegistrar{
		registerFn: func(_ context.Context, _, _ string) (domain.Profile, error) {
			return domain.Profile{}, nil
		},
	})

	req := httptest.NewRequest(http.MethodPost, "/profiles", bytes.NewBufferString(`{"name":"Alice"`))
	res := httptest.NewRecorder()

	h.ServeHTTP(res, req)

	if res.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusBadRequest)
	}
}

func TestProfileHandler_ServeHTTP_InvalidInput(t *testing.T) {
	t.Parallel()

	h := NewProfileHandler(stubProfileRegistrar{
		registerFn: func(_ context.Context, _, _ string) (domain.Profile, error) {
			return domain.Profile{}, profileregister.ErrInvalidInput
		},
	})

	req := httptest.NewRequest(http.MethodPost, "/profiles", bytes.NewBufferString(`{"name":" ","nickname":"ali"}`))
	res := httptest.NewRecorder()

	h.ServeHTTP(res, req)

	if res.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusBadRequest)
	}
}

func TestProfileHandler_ServeHTTP_InternalError(t *testing.T) {
	t.Parallel()

	h := NewProfileHandler(stubProfileRegistrar{
		registerFn: func(_ context.Context, _, _ string) (domain.Profile, error) {
			return domain.Profile{}, errors.New("failed")
		},
	})

	req := httptest.NewRequest(http.MethodPost, "/profiles", bytes.NewBufferString(`{"name":"Alice","nickname":"ali"}`))
	res := httptest.NewRecorder()

	h.ServeHTTP(res, req)

	if res.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusInternalServerError)
	}
}
