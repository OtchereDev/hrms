package errors

import "errors"

// Common application errors
var (
	ErrNotFound          = errors.New("resource not found")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrForbidden         = errors.New("forbidden")
	ErrBadRequest        = errors.New("bad request")
	ErrInternalServer    = errors.New("internal server error")
	ErrValidation        = errors.New("validation error")
	ErrDuplicate         = errors.New("duplicate resource")
	ErrInvalidInput      = errors.New("invalid input")
	ErrDatabaseOperation = errors.New("database operation failed")
)
