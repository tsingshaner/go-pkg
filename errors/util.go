package errors

import "errors"

// Code returns the code of the error. If the error does not have a code, it returns -1, false.
func Code[T Coder](e error) (code T, ok bool) {
	if e, ok := e.(interface{ Code() T }); ok {
		return e.Code(), true
	}

	return *new(T), false
}

func Status(e error) (status int, ok bool) {
	if e, ok := e.(interface{ Status() int }); ok {
		return e.Status(), true
	}

	return -1, false
}

// Extract try to extract the target error from the error chain, use errors.As, if the error is not the target, it returns nil.
func Extract[T error](e error) (target T, ok bool) {
	if e == nil {
		return target, false
	}

	if errors.As(e, &target) {
		return target, true
	}

	return target, false
}
