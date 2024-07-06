package errors

type basicError struct {
	msg  string
	code int
}

type BasicError interface {
	error
	Is(error) bool
	Code() int
}

var _ BasicError = &basicError{}

func (e *basicError) Error() string {
	return e.msg
}

func (e *basicError) Code() int {
	return e.code
}

func (e *basicError) Is(target error) bool {
	if targetErr, ok := target.(BasicError); ok {
		return targetErr.Code() == e.code
	}

	return false
}

func NewBasic(code int, msg string) error {
	return &basicError{msg, code}
}
