package errors

type UnauthorizedError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func NewUnauthorizedError(code, message string) *UnauthorizedError {
	return &UnauthorizedError{
		Code:    code,
		Message: message,
	}
}

func (e *UnauthorizedError) Error() string {
	return e.Message
}
