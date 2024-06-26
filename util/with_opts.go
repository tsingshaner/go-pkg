package util

type WithFn[T any] func(*T)

func ApplyOpts[T any](opts *T, fns ...WithFn[T]) {
	for _, fn := range fns {
		if fn != nil {
			fn(opts)
		}
	}
}

func BuildWithOpts[T any](opts *T, fns ...WithFn[T]) *T {
	ApplyOpts(opts, fns...)
	return opts
}
