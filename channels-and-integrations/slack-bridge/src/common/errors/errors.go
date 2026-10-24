package errors

import "errors"

// ErrValidation marks a request rejected by its Parse or Validate method.
var ErrValidation = errors.New("validation failed")
