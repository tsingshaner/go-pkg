package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type caseItem[T Coder] struct {
	err  error
	code T
	ok   bool
}

func TestCode(t *testing.T) {
	cases := []caseItem[int]{
		{errors.New("error"), 0, false},
		{NewBasic(1, "basic"), 1, true},
		{NewREST(1, 2, "rest"), 2, true},
	}

	for _, c := range cases {
		code, ok := Code[int](c.err)
		assert.Equal(t, c.code, code)
		assert.Equal(t, c.ok, ok)
	}

	cases2 := []caseItem[string]{
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

func TestStatus(t *testing.T) {
	restErr := NewREST(404, "x2", "rest")

	status, ok := Status(restErr)
	assert.Equal(t, 404, status)
	assert.True(t, ok)

	basicError := NewBasic(1, "basic")
	status, ok = Status(basicError)
	assert.Equal(t, -1, status)
	assert.False(t, ok)
}

func TestExtract(t *testing.T) {
	basicErr := NewBasic(1, "basic")
	restErr := errors.Join(NewREST(404, "x2", "rest"))
	err := errors.Join(restErr, basicErr)

	{
		e, ok := Extract[BasicError[int]](nil)
		assert.False(t, ok)
		assert.Nil(t, e)
	}

	{
		e, ok := Extract[BasicError[int]](errors.New("basic"))
		assert.False(t, ok)
		assert.Nil(t, e)
	}

	{
		e, ok := Extract[BasicError[string]](basicErr)
		assert.False(t, ok)
		assert.Nil(t, e)
	}

	{
		e, ok := Extract[BasicError[int]](err)
		assert.True(t, ok)
		assert.NotNil(t, e)
		assert.Equal(t, basicErr, e)
		assert.Equal(t, basicErr.Error(), e.Error())
		assert.Equal(t, 1, e.Code())
	}

	{
		e, ok := Extract[RESTError[string]](err)
		assert.True(t, ok)
		assert.NotNil(t, e)
		assert.True(t, errors.Is(restErr, e))
		assert.Equal(t, restErr.Error(), e.Error())
		assert.Equal(t, 404, e.Status())
	}
}
