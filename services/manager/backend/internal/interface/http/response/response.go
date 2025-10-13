// Package response contains helpers for shaping HTTP API responses.
package response

import "github.com/labstack/echo/v4"

// Envelope defines the standard API response structure.
type Envelope struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
	Error  *ErrorBody  `json:"error,omitempty"`
	Meta   interface{} `json:"meta,omitempty"`
}

// ErrorBody describes an error payload.
type ErrorBody struct {
	RequestID string        `json:"requestId"`
	Code      string        `json:"code,omitempty"`
	Message   string        `json:"message,omitempty"`
	Details   []ErrorDetail `json:"details,omitempty"`
}

// ErrorDetail provides granular validation error information.
type ErrorDetail struct {
	Field   string `json:"field,omitempty"`
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

// Success returns a success response with the provided HTTP status code.
func Success(c echo.Context, status int, data interface{}, meta interface{}) error {
	return c.JSON(status, Envelope{Status: "success", Data: data, Meta: meta})
}

// Failure returns a failure response with the provided details.
func Failure(c echo.Context, status int, requestID, code, message string, details []ErrorDetail) error {
	return c.JSON(status, Envelope{
		Status: "error",
		Error: &ErrorBody{
			RequestID: requestID,
			Code:      code,
			Message:   message,
			Details:   details,
		},
	})
}
