package errors

import "fmt"

// CustomError is an interface for custom application errors
type CustomError interface {
	Error() string
	Type() ErrorType
	Unwrap() error
}

// ErrorType defines categories of errors
type ErrorType string

const (
	TypeNotFound      ErrorType = "NOT_FOUND"
	TypeAlreadyExists ErrorType = "ALREADY_EXISTS"
	TypeValidation    ErrorType = "VALIDATION_ERROR"
	TypeInternal      ErrorType = "INTERNAL_ERROR"
	TypeUnauthorized  ErrorType = "UNAUTHORIZED"
	TypeForbidden     ErrorType = "FORBIDDEN"
	TypeConflict      ErrorType = "CONFLICT"
)

// appError is the concrete implementation of CustomError
type appError struct {
	err       error
	errorType ErrorType
}

// New create a new custom error with a type and an optional wrapped error
func New(t ErrorType, msg string, err error) CustomError {
	if err == nil {
		return &appError{err: fmt.Errorf("%s", msg), errorType: t}
	}
	return &appError{err: fmt.Errorf("%s: %w", msg, err), errorType: t}
}

// Error returns the error message
func (e *appError) Error() string {
	return e.err.Error()
}

// Type returns the type of the error
func (e *appError) Type() ErrorType {
	return e.errorType
}

// Unwrap returns the underlying wrapped error
func (e *appError) Unwrap() error {
	return e.err
}

// Convenience constructors for specific error types
func NewNotFoundError(msg string, err error) CustomError {
	return New(TypeNotFound, msg, err)
}

func NewAlreadyExistsError(msg string, err error) CustomError {
	return New(TypeAlreadyExists, msg, err)
}

func NewValidationError(msg string, err error) CustomError {
	return New(TypeValidation, msg, err)
}

func NewInternalError(msg string, err error) CustomError {
	return New(TypeInternal, msg, err)
}

func NewUnauthorizedError(msg string, err error) CustomError {
	return New(TypeUnauthorized, msg, err)
}

func NewForbiddenError(msg string, err error) CustomError {
	return New(TypeForbidden, msg, err)
}

func NewConflictError(msg string, err error) CustomError {
	return New(TypeConflict, msg, err)
}

// IsType checks if a given error is of a specific CustomError type
func IsType(err error, t ErrorType) bool {
	if customErr, ok := err.(CustomError); ok {
		return customErr.Type() == t
	}
	return false
}

// GetType retrieves the ErrorType from a CustomError, or return an empty string
func GetType(err error) ErrorType {
	if customErr, ok := err.(CustomError); ok {
		return customErr.Type()
	}
	return ""
}
