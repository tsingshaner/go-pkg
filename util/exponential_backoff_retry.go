package util

import (
	"time"
)

type RetryFunc[T any] func() (T, error)
type ExponentialOption struct {
	InitialDelay time.Duration
	Exponential  int
	Retry        int
}

// BackoffRetry 指数退避重试函数
// retry: 重试次数
// delay: 重试间隔
// f: 重试函数
func ExponentialBackoffRetry[T any](f RetryFunc[*T], opts *ExponentialOption) (*T, error) {
	var (
		err error
		res *T
	)

	for i := 0; i < opts.Retry; i++ {
		res, err = f()
		if err == nil {
			return res, nil
		}

		delay := opts.InitialDelay * time.Duration(i*opts.Exponential)
		// logger.Error(fmt.Sprintf("execute failed %v, will retry after %v", err, delay))
		time.Sleep(delay)
	}

	return res, err
}
