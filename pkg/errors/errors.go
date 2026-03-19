// Package errors provides error types for Higress Admin SDK.
package errors

import (
	"errors"
	"fmt"
)

// Common errors
var (
	// ErrNotFound indicates that the requested resource was not found
	ErrNotFound = errors.New("resource not found")
	// ErrConflict indicates a resource conflict (e.g., duplicate name)
	ErrConflict = errors.New("resource conflict")
	// ErrValidation indicates a validation error
	ErrValidation = errors.New("validation error")
	// ErrBusiness indicates a general business logic error
	ErrBusiness = errors.New("business error")
)

// ValidationError represents a validation error with details.
type ValidationError struct {
	Message string
	Field   string
}

// Error implements the error interface.
func (e *ValidationError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
	}
	return fmt.Sprintf("validation error: %s", e.Message)
}

// Is implements the errors.Is interface.
func (e *ValidationError) Is(target error) bool {
	return target == ErrValidation
}

// NewValidationError creates a new ValidationError.
func NewValidationError(message string) *ValidationError {
	return &ValidationError{Message: message}
}

// NewValidationErrorWithField creates a new ValidationError with a field name.
func NewValidationErrorWithField(message, field string) *ValidationError {
	return &ValidationError{Message: message, Field: field}
}

// NotFoundError represents an error when a resource is not found.
type NotFoundError struct {
	Resource string
	Name     string
}

// Error implements the error interface.
func (e *NotFoundError) Error() string {
	if e.Name != "" {
		return fmt.Sprintf("%s '%s' not found", e.Resource, e.Name)
	}
	return fmt.Sprintf("%s not found", e.Resource)
}

// Is implements the errors.Is interface.
func (e *NotFoundError) Is(target error) bool {
	return target == ErrNotFound
}

// NewNotFoundError creates a new NotFoundError.
func NewNotFoundError(resource, name string) *NotFoundError {
	return &NotFoundError{Resource: resource, Name: name}
}

// ResourceConflictError represents an error when a resource conflict occurs.
type ResourceConflictError struct {
	Resource string
	Message  string
}

// Error implements the error interface.
func (e *ResourceConflictError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("conflict on %s: %s", e.Resource, e.Message)
	}
	return fmt.Sprintf("conflict on %s", e.Resource)
}

// Is implements the errors.Is interface.
func (e *ResourceConflictError) Is(target error) bool {
	return target == ErrConflict
}

// NewResourceConflictError creates a new ResourceConflictError.
func NewResourceConflictError(resource, message string) *ResourceConflictError {
	return &ResourceConflictError{Resource: resource, Message: message}
}

// BusinessError represents a general business logic error.
type BusinessError struct {
	Message string
	Cause   error
}

// Error implements the error interface.
func (e *BusinessError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("business error: %s (cause: %v)", e.Message, e.Cause)
	}
	return fmt.Sprintf("business error: %s", e.Message)
}

// Is implements the errors.Is interface.
func (e *BusinessError) Is(target error) bool {
	return target == ErrBusiness
}

// Unwrap returns the underlying cause of the error.
func (e *BusinessError) Unwrap() error {
	return e.Cause
}

// NewBusinessError creates a new BusinessError.
func NewBusinessError(message string) *BusinessError {
	return &BusinessError{Message: message}
}

// NewBusinessErrorWithCause creates a new BusinessError with a cause.
func NewBusinessErrorWithCause(message string, cause error) *BusinessError {
	return &BusinessError{Message: message, Cause: cause}
}

// IsNotFound checks if the error is a NotFoundError.
func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}

// IsConflict checks if the error is a ResourceConflictError.
func IsConflict(err error) bool {
	return errors.Is(err, ErrConflict)
}

// IsValidation checks if the error is a ValidationError.
func IsValidation(err error) bool {
	return errors.Is(err, ErrValidation)
}

// IsBusiness checks if the error is a BusinessError.
func IsBusiness(err error) bool {
	return errors.Is(err, ErrBusiness)
}
