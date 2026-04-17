package requests

import "errors"

// ErrValidation wraps any error returned by a request's Validate method so
// callers can map it to an HTTP 400 via errors.Is.
var ErrValidation = errors.New("validation failed")
