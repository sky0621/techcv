package middleware

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/vibe-kanban/backend/internal/domain"
	"github.com/vibe-kanban/backend/internal/infrastructure/logger"
	"github.com/vibe-kanban/backend/internal/interface/http/response"
)

// ErrorHandler handles application errors and produces consistent API responses.
type ErrorHandler struct {
	logger *zap.Logger
}

// NewErrorHandler creates a new ErrorHandler instance.
func NewErrorHandler(logger *zap.Logger) *ErrorHandler {
	return &ErrorHandler{logger: logger}
}

// Handle implements echo's HTTPErrorHandler signature.
func (h *ErrorHandler) Handle(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	status := http.StatusInternalServerError
	code := "internal_error"
	message := http.StatusText(status)

	var appErr *domain.AppError
	if errors.As(err, &appErr) {
		status = appErr.StatusCode
		if status == 0 {
			status = http.StatusInternalServerError
		}
		if appErr.Code != "" {
			code = appErr.Code
		}
		if appErr.Message != "" {
			message = appErr.Message
		}
		err = appErr
	} else {
		var echoErr *echo.HTTPError
		if errors.As(err, &echoErr) {
			status = echoErr.Code
			if msg, ok := echoErr.Message.(string); ok && msg != "" {
				message = msg
			} else {
				message = http.StatusText(status)
			}
		} else {
			message = err.Error()
		}
	}

	requestID := c.Response().Header().Get(echo.HeaderXRequestID)
	log := logger.WithRequestID(h.logger, requestID)
	log.Error("request failed",
		zap.String("method", c.Request().Method),
		zap.String("path", c.Path()),
		zap.Int("status", status),
		zap.String("code", code),
		zap.Error(err),
	)

	if err := response.Failure(c, status, code, message); err != nil {
		log.Error("failed to send error response", zap.Error(err))
	}
}
