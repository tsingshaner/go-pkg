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
	return target.(BasicError).Code() == e.code
}

func NewBasic(code int, msg string) error {
	return &basicError{msg, code}
}
