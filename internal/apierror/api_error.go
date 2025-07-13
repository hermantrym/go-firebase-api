package apierror

import "net/http"

// APIError defines a standard error structure for our API responses.
type APIError struct {
	// Code is the HTTP status code. The `json:"-"` tag prevents it from being
	// rendered in the JSON response body.
	Code int `json:"-"`
	// Message is the user-friendly error message.
	Message string `json:"error"`
}

// Error implements the standard Go error interface, allowing APIError to be
// used as a regular error type.
func (e *APIError) Error() string {
	return e.Message
}

// NewAPIError creates a new instance of APIError.
func NewAPIError(code int, message string) *APIError {
	return &APIError{
		Code:    code,
		Message: message,
	}
}

// NewNotFoundError is a shortcut for creating a 404 Not Found error.
// It uses a default message if none is provided.
func NewNotFoundError(message string) *APIError {
	if message == "" {
		message = "The requested resource was not found"
	}

	return NewAPIError(http.StatusNotFound, message)
}

// NewInternalServerError is a shortcut for creating a 500 Internal Server Error.
// It uses a default message if none is provided.
func NewInternalServerError(message string) *APIError {
	if message == "" {
		message = "An unexpected internal error occurred"
	}

	return NewAPIError(http.StatusInternalServerError, message)
}

// NewBadRequestError is a shortcut for creating a 400 Bad Request error.
// It uses a default message if none is provided.
func NewBadRequestError(message string) *APIError {
	if message == "" {
		message = "Bad request"
	}

	return NewAPIError(http.StatusBadRequest, message)
}
