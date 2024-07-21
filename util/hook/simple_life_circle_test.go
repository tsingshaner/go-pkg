package hook

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleLifeCircle(t *testing.T) {
	type service struct {
		count int
	}

	slc := NewSimpleLifeCircle[*service]()
	slc.OnInit(func(s *service) error {
		s.count++
		return nil
	})

	slc.OnDone(func(s *service) error {
		s.count = 0
		return nil
	})

	s := &service{}

	//nolint:staticcheck // test pass a nil context
	errs := slc.Run(nil, s)
	assert.True(t, errors.Is(errs, ErrNilContext))

	ctx := context.TODO()
	errs = slc.Run(ctx, s)

	assert.Nil(t, errs)
	assert.Equal(t, 1, s.count)

	slc.Done()

	assert.Nil(t, <-slc.DoneError())
	assert.Equal(t, 0, s.count)

	assert.True(t, errors.Is(slc.Run(ctx, s), ErrLifeCircleIsDead))
}

func TestSimpleLifeCircleWithInitError(t *testing.T) {

	err := errors.New("init error")

	slc := NewSimpleLifeCircle[any]()
	slc.OnInit(func(_ any) error {
		return err
	})

	ctx := context.TODO()
	errs := slc.Run(ctx, struct{}{})

	assert.NotNil(t, errs)
	assert.True(t, errors.Is(errs, err))
}
