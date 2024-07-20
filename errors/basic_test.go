package errors

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasicError(t *testing.T) {
	baseErr := NewBasic(1, "base error")

	t.Run("error", func(t *testing.T) {
		assert.Implements(t, (*error)(nil), baseErr)
		assert.Implements(t, (*BasicError[int])(nil), baseErr)
		assert.Equal(t, baseErr.Error(), "base error")
		assert.Equal(t, baseErr.(BasicError[int]).Code(), 1)
	})

	t.Run("is", func(t *testing.T) {
		assert.True(t, errors.Is(baseErr, baseErr))
		assert.False(t, errors.Is(baseErr, NewBasic(1, "base error")))
		assert.False(t, errors.Is(baseErr, NewBasic(2, "base error")))
		assert.False(t, errors.Is(baseErr, errors.New("base error")))
	})

	t.Run("as", func(t *testing.T) {
		var target BasicError[int]
		assert.True(t, errors.As(baseErr, &target))
		assert.Equal(t, target.Code(), 1)
	})

	t.Run("wrap", func(t *testing.T) {
		wrapErr := fmt.Errorf("wrapped error: %w", baseErr)
		assert.True(t, errors.Is(wrapErr, baseErr))
		assert.Equal(t, wrapErr.Error(), "wrapped error: base error")

		var target BasicError[int]
		assert.True(t, errors.As(wrapErr, &target))
		assert.Equal(t, target.Code(), 1)
	})

	t.Run("wrap/nest", func(t *testing.T) {
		wrapErr1 := fmt.Errorf("wrapped error1: %w", baseErr)
		wrapErr2 := fmt.Errorf("wrapped error2: %w", wrapErr1)
		assert.True(t, errors.Is(wrapErr2, baseErr))
		assert.True(t, errors.Is(wrapErr2, wrapErr1))
		assert.Equal(t, wrapErr2.Error(), "wrapped error2: wrapped error1: base error")

		var target BasicError[int]
		assert.True(t, errors.As(wrapErr2, &target))
		assert.Equal(t, target.Code(), 1)
	})

	t.Run("wrap/slices", func(t *testing.T) {
		originErr := errors.New("origin error")
		wrapErr1 := fmt.Errorf("wrapped origin error: %w", originErr)

		wrapErrs := fmt.Errorf("wrapped error3: (%w, %w)", wrapErr1, baseErr)
		assert.True(t, errors.Is(wrapErrs, baseErr))
	})
}
