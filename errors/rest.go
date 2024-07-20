package errors

type restError[T Coder] struct {
	basic  *basicError[T]
	status int
}

type RESTError[T Coder] interface {
	BasicError[T]
	Status() int
}

var _ RESTError[int] = &restError[int]{}

func (e *restError[T]) Code() T {
	return e.basic.Code()
}

func (e *restError[T]) Error() string {
	return e.basic.Error()
}

func (e *restError[T]) Status() int {
	return e.status
}

func (e *restError[T]) Is(target error) bool {
	if restErr, ok := target.(RESTError[T]); ok {
		return e == restErr
	}

	return false
}

func NewREST[T Coder](status int, code T, msg string) error {
	return &restError[T]{&basicError[T]{msg, code}, status}
}
