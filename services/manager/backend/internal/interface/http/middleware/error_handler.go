package middleware

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/sky0621/techcv/manager/backend/internal/domain"
	"github.com/sky0621/techcv/manager/backend/internal/infrastructure/logger"
	"github.com/sky0621/techcv/manager/backend/internal/interface/http/response"
)

// ErrorHandler handles application errors and produces consistent API responses.
type ErrorHandler struct {
	logger *slog.Logger
}

// NewErrorHandler creates a new ErrorHandler instance.
func NewErrorHandler(logger *slog.Logger) *ErrorHandler {
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
		slog.String("method", c.Request().Method),
		slog.String("path", c.Path()),
		slog.Int("status", status),
		slog.String("code", code),
		slog.Any("error", err),
	)

	var details []response.ErrorDetail
	if appErr != nil {
		for _, d := range appErr.Details {
			details = append(details, response.ErrorDetail{
				Field:   d.Field,
				Code:    d.Code,
				Message: d.Message,
			})
		}
	}

	if err := response.Failure(c, status, requestID, code, message, details); err != nil {
		log.Error("failed to send error response", slog.Any("error", err))
	}
}
