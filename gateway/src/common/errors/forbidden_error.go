package errors

import "fmt"

// ForbiddenError maps to HTTP 403.
type ForbiddenError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func NewForbiddenError(code, message string) *ForbiddenError {
	return &ForbiddenError{Code: code, Message: message}
}

func (e *ForbiddenError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}
