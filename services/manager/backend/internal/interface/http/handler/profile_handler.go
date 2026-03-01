package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/sky0621/techcv/services/manager/backend/internal/domain/model"
	profileregister "github.com/sky0621/techcv/services/manager/backend/internal/usecase/profile/register"
)

type profileRegistrar interface {
	Register(ctx context.Context, name, nickname string) (model.Profile, error)
}

type ProfileHandler struct {
	profileService profileRegistrar
}

func NewProfileHandler(profileService profileRegistrar) ProfileHandler {
	return ProfileHandler{profileService: profileService}
}

type registerProfileRequest struct {
	Name     string `json:"name"`
	Nickname string `json:"nickname"`
}

func (h ProfileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req registerProfileRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := dec.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	profile, err := h.profileService.Register(r.Context(), req.Name, req.Nickname)
	if err != nil {
		if errors.Is(err, profileregister.ErrInvalidInput) {
			writeError(w, http.StatusBadRequest, "name and nickname are required")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to register profile")
		return
	}

	writeJSON(w, http.StatusCreated, profile)
}

func writeError(w http.ResponseWriter, statusCode int, message string) {
	writeJSON(w, statusCode, map[string]string{"error": message})
}

func writeJSON(w http.ResponseWriter, statusCode int, data any) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(data); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"failed to encode response"}` + "\n"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, _ = w.Write(buf.Bytes())
}
