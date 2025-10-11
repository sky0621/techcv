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
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Success returns a success response with the provided HTTP status code.
func Success(c echo.Context, status int, data interface{}) error {
	return c.JSON(status, Envelope{Status: "success", Data: data})
}

// Failure returns a failure response with the provided details.
func Failure(c echo.Context, status int, code, message string) error {
	return c.JSON(status, Envelope{
		Status: "error",
		Error: &ErrorBody{
			Code:    code,
			Message: message,
		},
	})
}
