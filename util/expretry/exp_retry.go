package exp_retry

import (
	"math"
	"time"

	"github.com/tsingshaner/go-pkg/log/console"
	"github.com/tsingshaner/go-pkg/util"
)

type Runner[T any] func() (T, error)

type ExpRetrySvc[T any] struct {
	Runner       Runner[T]
	Delay        time.Duration
	RetryTimes   int
	ErrorHandler func(e error, times int)
}

func New[T any](runner Runner[T], fns ...util.WithFn[ExpRetrySvc[T]]) *ExpRetrySvc[T] {
	return util.BuildWithOpts(&ExpRetrySvc[T]{
		Runner:     runner,
		Delay:      time.Second,
		RetryTimes: 3,
		ErrorHandler: func(e error, times int) {
			console.Error("Attempt %d failed: %w", times, e)
		},
	}, fns...)
}

func (e *ExpRetrySvc[T]) Run() (res T, err error) {
	for i := 0; i < e.RetryTimes; i++ {
		if res, err = e.Runner(); err == nil {
			return res, nil
		}

		e.ErrorHandler(err, i)
		time.Sleep(e.Delay * time.Duration(math.Pow(2, float64(i))))
	}

	return res, err
}
