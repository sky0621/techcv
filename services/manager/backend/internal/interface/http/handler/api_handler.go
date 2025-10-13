// Package handler exposes HTTP handlers that satisfy the OpenAPI contract.
package handler

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/sky0621/techcv/manager/backend/internal/domain"
	openapi "github.com/sky0621/techcv/manager/backend/internal/interface/http/openapi"
	"github.com/sky0621/techcv/manager/backend/internal/interface/http/response"
	"github.com/sky0621/techcv/manager/backend/internal/usecase/auth"
)

// HealthUsecase defines the behavior required by the handler.
type HealthUsecase interface {
	Check(ctx context.Context) (*domain.HealthStatus, error)
}

// RegisterUsecase defines the registration usecase contract.
type RegisterUsecase interface {
	Execute(ctx context.Context, in auth.RegisterInput) (auth.RegisterOutput, error)
}

// VerifyUsecase defines the email verification contract.
type VerifyUsecase interface {
	Execute(ctx context.Context, in auth.VerifyInput) (auth.VerifyOutput, error)
}

// Handler implements the OpenAPI server interface.
type Handler struct {
	health   HealthUsecase
	register RegisterUsecase
	verify   VerifyUsecase
}

// NewHandler creates a new API handler instance.
func NewHandler(health HealthUsecase, register RegisterUsecase, verify VerifyUsecase) *Handler {
	return &Handler{
		health:   health,
		register: register,
		verify:   verify,
	}
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

// PostAuthRegister handles guest registration.
func (h *Handler) PostAuthRegister(c echo.Context) error {
	var req openapi.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return domain.NewValidation(domain.ErrorCodeInvalidRequest, "リクエスト形式が正しくありません").WithDetails(
			domain.ErrorDetail{Field: "body", Code: domain.ErrorCodeInvalidJSON, Message: "JSONの解析に失敗しました"},
		)
	}

	out, err := h.register.Execute(c.Request().Context(), auth.RegisterInput{
		Email:                req.Email,
		Password:             req.Password,
		PasswordConfirmation: req.PasswordConfirmation,
	})
	if err != nil {
		return err
	}

	data := map[string]interface{}{
		"message":    out.Message,
		"expires_at": out.ExpiresAt,
	}

	meta := map[string]interface{}{
		"requestId": c.Response().Header().Get(echo.HeaderXRequestID),
	}

	return response.Success(c, http.StatusOK, data, meta)
}

// PostAuthVerify finalizes registration and authenticates the user.
func (h *Handler) PostAuthVerify(c echo.Context) error {
	var req openapi.VerifyRequest
	if err := c.Bind(&req); err != nil {
		return domain.NewValidation(domain.ErrorCodeInvalidRequest, "リクエスト形式が正しくありません").WithDetails(
			domain.ErrorDetail{Field: "body", Code: domain.ErrorCodeInvalidJSON, Message: "JSONの解析に失敗しました"},
		)
	}

	out, err := h.verify.Execute(c.Request().Context(), auth.VerifyInput{Token: req.Token})
	if err != nil {
		return err
	}

	user := out.User
	payload := map[string]interface{}{
		"message":    out.Message,
		"auth_token": out.AuthToken,
		"user": map[string]interface{}{
			"id":                user.ID,
			"email":             user.Email,
			"name":              user.Name,
			"bio":               user.Bio,
			"is_active":         user.IsActive,
			"email_verified_at": user.EmailVerifiedAt,
			"last_login_at":     user.LastLoginAt,
			"created_at":        user.CreatedAt,
			"updated_at":        user.UpdatedAt,
		},
	}

	meta := map[string]interface{}{
		"requestId": c.Response().Header().Get(echo.HeaderXRequestID),
	}

	return response.Success(c, http.StatusOK, payload, meta)
}
