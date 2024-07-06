package errors

type restError struct {
	basic  *basicError
	status int
}

type RESTError interface {
	BasicError
	Status() int
}

var _ RESTError = &restError{}

func (e *restError) Code() int {
	return e.basic.Code()
}

func (e *restError) Error() string {
	return e.basic.Error()
}

func (e *restError) Status() int {
	return e.status
}

func (e *restError) Is(target error) bool {
	if restErr, ok := target.(RESTError); ok && restErr.Status() == e.status {
		return e.basic.code == restErr.Code()
	}

	return false
}

func NewREST(status, code int, msg string) error {
	return &restError{&basicError{msg, code}, status}
}
