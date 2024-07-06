package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCode(t *testing.T) {
	cases := []struct {
		err  error
		code int
		ok   bool
	}{
		{errors.New("error"), -1, false},
		{NewBasic(1, "basic"), 1, true},
		{NewREST(1, 2, "rest"), 2, true},
	}

	for _, c := range cases {
		code, ok := Code(c.err)
		assert.Equal(t, c.code, code)
		assert.Equal(t, c.ok, ok)
	}
}
