package errors

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRESTError(t *testing.T) {
	restErr := NewREST(http.StatusInternalServerError, 10001, "rest")

	t.Run("error", func(t *testing.T) {
		assert.Implements(t, (*error)(nil), restErr)
		assert.Implements(t, (*RESTError)(nil), restErr)
		assert.Equal(t, restErr.Error(), "rest")
		assert.Equal(t, restErr.(RESTError).Code(), 10001)
		assert.Equal(t, restErr.(RESTError).Status(), http.StatusInternalServerError)
	})

	t.Run("is", func(t *testing.T) {
		assert.True(t, errors.Is(
			restErr,
			NewREST(http.StatusInternalServerError, 10001, "rest"),
		))

		assert.False(t, errors.Is(
			restErr,
			NewREST(http.StatusBadRequest, 10001, "rest"),
		))

		assert.False(t, errors.Is(
			restErr,
			NewBasic(10001, "rest"),
		))
	})
}
