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
		{errors.New("error"), 0, false},
		{NewBasic(1, "basic"), 1, true},
		{NewREST(1, 2, "rest"), 2, true},
	}

	for _, c := range cases {
		code, ok := Code[int](c.err)
		assert.Equal(t, c.code, code)
		assert.Equal(t, c.ok, ok)
	}

	cases2 := []struct {
		err  error
		code string
		ok   bool
	}{
		{errors.New("error"), "", false},
		{NewBasic("x1", "basic"), "x1", true},
		{NewREST(404, "x2", "rest"), "x2", true},
	}

	for _, c := range cases2 {
		code, ok := Code[string](c.err)
		assert.Equal(t, c.code, code)
		assert.Equal(t, c.ok, ok)
	}
}
