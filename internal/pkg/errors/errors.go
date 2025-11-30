package errors

import "errors"

// Common reusable errors for adapters and use cases.
var (
	ErrNotFound       = errors.New("not found")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrNotImplemented = errors.New("not implemented")
)
