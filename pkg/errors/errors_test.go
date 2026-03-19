package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestValidationError 测试 ValidationError
// 运行命令: go test -v -run TestValidationError ./pkg/errors/
func TestValidationError(t *testing.T) {
	t.Run("create validation error", func(t *testing.T) {
		err := NewValidationError("field is required")
		assert.Equal(t, "validation error: field is required", err.Error())
		assert.True(t, IsValidation(err))
	})

	t.Run("validation error with field", func(t *testing.T) {
		err := NewValidationErrorWithField("is required", "name")
		assert.Equal(t, "validation error on field 'name': is required", err.Error())
		assert.True(t, IsValidation(err))
	})

	t.Run("errors.Is support", func(t *testing.T) {
		err := NewValidationError("test")
		assert.True(t, errors.Is(err, ErrValidation))
	})
}

// TestNotFoundError 测试 NotFoundError
// 运行命令: go test -v -run TestNotFoundError ./pkg/errors/
func TestNotFoundError(t *testing.T) {
	t.Run("create not found error", func(t *testing.T) {
		err := NewNotFoundError("Domain", "example.com")
		assert.Equal(t, "Domain 'example.com' not found", err.Error())
		assert.True(t, IsNotFound(err))
	})

	t.Run("not found error without name", func(t *testing.T) {
		err := &NotFoundError{Resource: "Domain"}
		assert.Equal(t, "Domain not found", err.Error())
	})

	t.Run("errors.Is support", func(t *testing.T) {
		err := NewNotFoundError("Domain", "test")
		assert.True(t, errors.Is(err, ErrNotFound))
	})
}

// TestResourceConflictError 测试 ResourceConflictError
// 运行命令: go test -v -run TestResourceConflictError ./pkg/errors/
func TestResourceConflictError(t *testing.T) {
	t.Run("create conflict error", func(t *testing.T) {
		err := NewResourceConflictError("Domain", "already exists")
		assert.Equal(t, "conflict on Domain: already exists", err.Error())
		assert.True(t, IsConflict(err))
	})

	t.Run("conflict error without message", func(t *testing.T) {
		err := &ResourceConflictError{Resource: "Domain"}
		assert.Equal(t, "conflict on Domain", err.Error())
	})

	t.Run("errors.Is support", func(t *testing.T) {
		err := NewResourceConflictError("Domain", "test")
		assert.True(t, errors.Is(err, ErrConflict))
	})
}

// TestBusinessError 测试 BusinessError
// 运行命令: go test -v -run TestBusinessError ./pkg/errors/
func TestBusinessError(t *testing.T) {
	t.Run("create business error", func(t *testing.T) {
		err := NewBusinessError("operation failed")
		assert.Equal(t, "business error: operation failed", err.Error())
		assert.True(t, IsBusiness(err))
	})

	t.Run("business error with cause", func(t *testing.T) {
		cause := errors.New("underlying error")
		err := NewBusinessErrorWithCause("operation failed", cause)
		assert.Contains(t, err.Error(), "operation failed")
		assert.Contains(t, err.Error(), "underlying error")
		assert.True(t, IsBusiness(err))
	})

	t.Run("errors.Is support", func(t *testing.T) {
		err := NewBusinessError("test")
		assert.True(t, errors.Is(err, ErrBusiness))
	})

	t.Run("Unwrap support", func(t *testing.T) {
		cause := errors.New("underlying error")
		err := NewBusinessErrorWithCause("test", cause)
		unwrapped := errors.Unwrap(err)
		assert.Equal(t, cause, unwrapped)
	})
}

// TestErrorTypeChecks 测试错误类型检查函数
// 运行命令: go test -v -run TestErrorTypeChecks ./pkg/errors/
func TestErrorTypeChecks(t *testing.T) {
	t.Run("IsNotFound", func(t *testing.T) {
		assert.True(t, IsNotFound(NewNotFoundError("test", "test")))
		assert.False(t, IsNotFound(NewValidationError("test")))
	})

	t.Run("IsConflict", func(t *testing.T) {
		assert.True(t, IsConflict(NewResourceConflictError("test", "test")))
		assert.False(t, IsConflict(NewValidationError("test")))
	})

	t.Run("IsValidation", func(t *testing.T) {
		assert.True(t, IsValidation(NewValidationError("test")))
		assert.False(t, IsValidation(NewNotFoundError("test", "test")))
	})

	t.Run("IsBusiness", func(t *testing.T) {
		assert.True(t, IsBusiness(NewBusinessError("test")))
		assert.False(t, IsBusiness(NewValidationError("test")))
	})
}
