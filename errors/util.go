package errors

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

func Msg(e error) (msg string, ok bool) {
	if e, ok := e.(interface{ Msg() string }); ok {
		return e.Msg(), true
	}

	return "", false
}
