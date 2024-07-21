package hook

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHook(t *testing.T) {
	count := 0
	h := New[any]()

	h.On(func(_ any) error {
		count++
		return nil
	})

	h.On(func(_ any) error {
		count += 2
		return nil
	})

	_ = h.Trigger(struct{}{})
	h.Clear()

	assert.Equal(t, 3, count)
}

func TestHookCancel(t *testing.T) {
	count := 0
	h := New[any]()

	off1 := h.On(func(_ any) error {
		count++
		return nil
	})

	off2 := h.On(func(_ any) error {
		count += 2
		return nil
	})

	_ = h.Trigger(struct{}{})
	assert.Equal(t, 3, count)

	off2()
	count = 0
	_ = h.Trigger(struct{}{})
	assert.Equal(t, 1, count)

	off1()
	count = 0
	_ = h.Trigger(struct{}{})
	assert.Equal(t, 0, count)
}

func TestHookWithData(t *testing.T) {
	h := New[*int]()

	h.On(func(data *int) error {
		*data += 1
		return nil
	})

	count := new(int)
	_ = h.Trigger(count)
	h.Clear()

	assert.Equal(t, 1, *count)
}

func TestHookWithError(t *testing.T) {
	h := New[any]()
	err := errors.New("error1")

	h.On(func(_ any) error {
		return err
	})

	h.On(func(_ any) error { return nil })

	h.On(func(_ any) error {
		return errors.New("error2")
	})

	errs := h.Trigger(struct{}{})
	h.Clear()

	assert.Error(t, errs)
	assert.Contains(t, errs.Error(), "error1")
	assert.Contains(t, errs.Error(), "error2")
	assert.True(t, errors.Is(errs, err))
}
