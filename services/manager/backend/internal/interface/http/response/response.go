// Package response contains helpers for shaping HTTP API responses.
package response

import "github.com/labstack/echo/v4"

// ErrorBody describes an error payload.
type ErrorBody struct {
	RequestID string        `json:"requestId"`
	Code      string        `json:"code,omitempty"`
	Details   []ErrorDetail `json:"details,omitempty"`
}

// ErrorDetail provides granular validation error information.
type ErrorDetail struct {
	Field   string `json:"field,omitempty"`
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

// Success returns a success response with the provided HTTP status code.
func Success(c echo.Context, status int, data interface{}) error {
	return c.JSON(status, data)
}

// Failure returns a failure response with the provided details.
func Failure(c echo.Context, status int, body ErrorBody) error {
	return c.JSON(status, body)
}
