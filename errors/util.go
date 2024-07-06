package errors

// Code returns the code of the error. If the error does not have a code, it returns -1, false.
func Code(e error) (code int, ok bool) {
	if e, ok := e.(interface{ Code() int }); ok {
		return e.Code(), true
	}

	return -1, false
}
