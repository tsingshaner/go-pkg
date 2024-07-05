package errors

type restError struct {
	msg    string
	status string
	basic  BasicError
}

type RestError interface {
	BasicError
	Status() string
}

var _ RestError = &restError{}

func (e *restError) Code() int {
	return e.basic.Code()
}

func (e *restError) Error() string {
	return e.msg
}

func (e *restError) Status() string {
	return e.status
}

func (e *restError) Is(target error) bool {
	return e.msg == target.Error() && target.(RestError).Status() == e.status
}

func (e *restError) Unwrap() error {
	return e.basic
}

func NewRestError(status, msg string, basic BasicError) error {
	return &restError{msg, status, basic}
}
