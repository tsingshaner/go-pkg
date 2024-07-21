package hook

import (
	"context"
	"errors"
)

type SimpleLifeCircle[T any] interface {
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

type simpleLifeCircle[T any] struct {
	cancel     func()
	dead       bool
	init       Hooker[T]
	done       Hooker[T]
	doneErrors chan error
}

func NewSimpleLifeCircle[T any]() SimpleLifeCircle[T] {
	return &simpleLifeCircle[T]{
		init:       New[T](),
		done:       New[T](),
		doneErrors: make(chan error, 1),
	}
}

func (s *simpleLifeCircle[T]) OnInit(handler Handler[T]) (off func()) {
	return s.init.On(handler)
}

func (s *simpleLifeCircle[T]) OnDone(handler Handler[T]) (off func()) {
	return s.done.On(handler)
}

var (
	ErrNilContext       = errors.New("SimpleLifeCircle.Run with a nil context.Context")
	ErrLifeCircleIsDead = errors.New("SimpleLifeCircle.Run with a dead life circle")
)

func (s *simpleLifeCircle[T]) Run(ctx context.Context, data T) error {
	if ctx == nil {
		return ErrNilContext
	} else if s.dead {
		return ErrLifeCircleIsDead
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

func (s *simpleLifeCircle[T]) Done() {
	if s.cancel != nil && !s.dead {
		s.cancel()
	}
}

func (s *simpleLifeCircle[T]) DoneError() <-chan error {
	return s.doneErrors
}
