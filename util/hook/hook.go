package hook

import "errors"

type Handler[T any] func(T) error

type Hooker[T any] interface {
	On(Handler[T]) (off func())
	Trigger(T) error
	Clear()
}

type hook[T any] struct {
	handlers map[int]Handler[T]
	nextId   int
}

func New[T any]() Hooker[T] {
	return &hook[T]{
		handlers: make(map[int]Handler[T]),
	}
}

func (h *hook[T]) Clear() {
	h.handlers = make(map[int]Handler[T])
}

func (h *hook[T]) On(handler Handler[T]) (off func()) {
	id := h.nextId
	h.handlers[id] = handler
	h.nextId++
	return func() {
		h.off(id)
	}
}

func (h *hook[T]) Trigger(data T) error {
	var errs error

	for _, handler := range h.handlers {
		errs = errors.Join(errs, handler(data))
	}

	return errs
}

func (h *hook[T]) off(id int) {
	delete(h.handlers, id)
}
