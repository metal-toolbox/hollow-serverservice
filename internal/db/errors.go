package db

import (
	"errors"
	"fmt"
)

var (
	// ErrNotFound is the error returned when a record is not found in the store
	ErrNotFound = errors.New("record not found in datastore")
)

// NotFoundError is the error returned when a record is not found in the datastore
type NotFoundError struct {
	Message string
}

func newNotFoundError(notFound string) *NotFoundError {
	return &NotFoundError{
		Message: fmt.Sprintf("%s not found", notFound),
	}
}

// Error returns a string showing the URI and the status code
func (e *NotFoundError) Error() string {
	return e.Message
}

// ValidationError is returned when an object fails Validation
type ValidationError struct {
	Message string
}

// Error returns a string showing the URI and the status code
func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed: %s", e.Message)
}

func newValidationError(msg string) *ValidationError {
	return &ValidationError{
		Message: msg,
	}
}

func requiredFieldMissing(model, field string) *ValidationError {
	return newValidationError(
		fmt.Sprintf("%s is a required %s attribute", field, model),
	)
}
