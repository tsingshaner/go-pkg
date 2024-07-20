package errors

import "golang.org/x/exp/constraints"

type Coder interface {
	constraints.Integer | string
}

type basicError[T Coder] struct {
	msg  string
	code T
}

type BasicError[T Coder] interface {
	error
	Is(error) bool
	Code() T
}

var _ BasicError[string] = &basicError[string]{"", ""}

func (e *basicError[T]) Error() string {
	return e.msg
}

func (e *basicError[T]) Code() T {
	return e.code
}

func (e *basicError[T]) Is(target error) bool {
	if targetErr, ok := target.(BasicError[T]); ok {
		return targetErr == e
	}

	return false
}

func NewBasic[T Coder](code T, msg string) error {
	return &basicError[T]{msg, code}
}
