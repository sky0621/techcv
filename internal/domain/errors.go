package domain

import (
	"errors"
	"net/http"
)

// AppError captures domain specific error details that can be surfaced over HTTP.
type AppError struct {
	Code       string
	Message    string
	StatusCode int
	Err        error
}

// Error implements the error interface.
func (e *AppError) Error() string {
	if e == nil {
		return ""
	}
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

// Unwrap exposes the underlying error.
func (e *AppError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

// NewNotFound returns a new not found error.
func NewNotFound(code, message string) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: http.StatusNotFound,
	}
}

// NewValidation returns a new validation error.
func NewValidation(code, message string) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: http.StatusBadRequest,
	}
}

// NewInternal returns a new internal server error.
func NewInternal(code, message string, err error) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: http.StatusInternalServerError,
		Err:        err,
	}
}

// IsAppError tests whether an error is an AppError.
func IsAppError(err error) bool {
	var target *AppError
	return errors.As(err, &target)
}
