package apperrors

import "errors"

var (
	ErrRecordNotFound = errors.New("not found")
	ErrInvalidFields  = errors.New("validation error")
)
