package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sky0621/techcv/services/manager/backend/internal/domain/model"
	profileregister "github.com/sky0621/techcv/services/manager/backend/internal/usecase/profile/register"
)

type mockProfileRegistrar struct {
	registerFn func(ctx context.Context, name, nickname string) (model.Profile, error)
	callCount  int
	gotName    string
	gotNick    string
}

func (m *mockProfileRegistrar) Register(ctx context.Context, name, nickname string) (model.Profile, error) {
	m.callCount++
	m.gotName = name
	m.gotNick = nickname
	return m.registerFn(ctx, name, nickname)
}

func TestProfileHandler_ServeHTTP_Created(t *testing.T) {
	t.Parallel()

	mockRegistrar := &mockProfileRegistrar{
		registerFn: func(_ context.Context, name, nickname string) (model.Profile, error) {
			return model.Profile{
				ID:       "profile_1",
				Name:     name,
				Nickname: nickname,
			}, nil
		},
	}
	h := NewProfileHandler(mockRegistrar)

	req := httptest.NewRequest(http.MethodPost, "/profiles", bytes.NewBufferString(`{"name":"Alice","nickname":"ali"}`))
	res := httptest.NewRecorder()

	h.ServeHTTP(res, req)

	if res.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusCreated)
	}
	var body model.Profile
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Name != "Alice" || body.Nickname != "ali" {
		t.Fatalf("unexpected response: %+v", body)
	}
	if mockRegistrar.callCount != 1 {
		t.Fatalf("register call count = %d, want 1", mockRegistrar.callCount)
	}
	if mockRegistrar.gotName != "Alice" || mockRegistrar.gotNick != "ali" {
		t.Fatalf("register args = (%q, %q), want (%q, %q)", mockRegistrar.gotName, mockRegistrar.gotNick, "Alice", "ali")
	}
}

func TestProfileHandler_ServeHTTP_BadRequest(t *testing.T) {
	t.Parallel()

	mockRegistrar := &mockProfileRegistrar{
		registerFn: func(_ context.Context, _, _ string) (model.Profile, error) {
			return model.Profile{}, nil
		},
	}
	h := NewProfileHandler(mockRegistrar)

	req := httptest.NewRequest(http.MethodPost, "/profiles", bytes.NewBufferString(`{"name":"Alice"`))
	res := httptest.NewRecorder()

	h.ServeHTTP(res, req)

	if res.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusBadRequest)
	}
	if mockRegistrar.callCount != 0 {
		t.Fatalf("register call count = %d, want 0", mockRegistrar.callCount)
	}
}

func TestProfileHandler_ServeHTTP_InvalidInput(t *testing.T) {
	t.Parallel()

	mockRegistrar := &mockProfileRegistrar{
		registerFn: func(_ context.Context, _, _ string) (model.Profile, error) {
			return model.Profile{}, profileregister.ErrInvalidInput
		},
	}
	h := NewProfileHandler(mockRegistrar)

	req := httptest.NewRequest(http.MethodPost, "/profiles", bytes.NewBufferString(`{"name":" ","nickname":"ali"}`))
	res := httptest.NewRecorder()

	h.ServeHTTP(res, req)

	if res.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusBadRequest)
	}
	if mockRegistrar.callCount != 1 {
		t.Fatalf("register call count = %d, want 1", mockRegistrar.callCount)
	}
}

func TestProfileHandler_ServeHTTP_InternalError(t *testing.T) {
	t.Parallel()

	mockRegistrar := &mockProfileRegistrar{
		registerFn: func(_ context.Context, _, _ string) (model.Profile, error) {
			return model.Profile{}, errors.New("failed")
		},
	}
	h := NewProfileHandler(mockRegistrar)

	req := httptest.NewRequest(http.MethodPost, "/profiles", bytes.NewBufferString(`{"name":"Alice","nickname":"ali"}`))
	res := httptest.NewRecorder()

	h.ServeHTTP(res, req)

	if res.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusInternalServerError)
	}
	if mockRegistrar.callCount != 1 {
		t.Fatalf("register call count = %d, want 1", mockRegistrar.callCount)
	}
}
