package hook

import (
	"context"
	"errors"
)

type SimpleLifecycle[T any] interface {
	// Run starts the life circle with the given data, if context is nil, it will return an error.
	//
	// Will emit Init event when call Run
	// Run again will return an error if the life circle is dead.
	Run(context.Context, T) error
	// OnInit registers a handler to be called when the life circle starts.
	OnInit(Handler[T]) (off func())
	// OnDone registers a handler to be called when the life circle ends.
	OnDone(Handler[T]) (off func())
	// OnError returns a channel that receives an error when the life circle ends.
	DoneError() <-chan error
	// Done end the life circle.
	Done()
}

type simpleLifecycle[T any] struct {
	cancel     func()
	dead       bool
	init       Hooker[T]
	done       Hooker[T]
	doneErrors chan error
}

func NewSimpleLifecycle[T any]() SimpleLifecycle[T] {
	return &simpleLifecycle[T]{
		init:       New[T](),
		done:       New[T](),
		doneErrors: make(chan error, 1),
	}
}

func (s *simpleLifecycle[T]) OnInit(handler Handler[T]) (off func()) {
	return s.init.On(handler)
}

func (s *simpleLifecycle[T]) OnDone(handler Handler[T]) (off func()) {
	return s.done.On(handler)
}

var (
	ErrNilContext      = errors.New("SimpleLifecycle.Run with a nil context.Context")
	ErrLifecycleIsDead = errors.New("SimpleLifecycle.Run with a dead life circle")
)

func (s *simpleLifecycle[T]) Run(ctx context.Context, data T) error {
	if ctx == nil {
		return ErrNilContext
	} else if s.dead {
		return ErrLifecycleIsDead
	}

	if err := s.init.Trigger(data); err != nil {
		return err
	}

	ctx, s.cancel = context.WithCancel(ctx)

	go func() {
		<-ctx.Done()
		s.dead = true
		s.doneErrors <- s.done.Trigger(data)
		close(s.doneErrors)
	}()

	return nil
}

func (s *simpleLifecycle[T]) Done() {
	if s.cancel != nil && !s.dead {
		s.cancel()
	}
}

func (s *simpleLifecycle[T]) DoneError() <-chan error {
	return s.doneErrors
}
