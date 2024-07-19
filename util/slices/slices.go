package slices

func Filter[T any](s []T, f func(T) bool) []T {
	r := make([]T, 0, len(s))
	for _, v := range s {
		if f(v) {
			r = append(r, v)
		}
	}

	return r
}

func Map[T, R any](s []T, f func(T) R) []R {
	r := make([]R, 0, len(s))
	for _, v := range s {
		r = append(r, f(v))
	}

	return r
}

// LastIndex returns the index of the last occurrence of the specified value in the slice,
// or -1 if the value is not present in the slice.
func LastIndex[T comparable](s []T, want T) int {
	i := len(s) - 1
	for ; i >= 0; i-- {
		if s[i] == want {
			return i
		}
	}

	return i
}

// LastIndex returns the index of the last occurrence of the specified value in the slice,
// or -1 if the value is not present in the slice.
func LastIndexFunc[T any](s []T, checker func(item T) bool) int {
	i := len(s) - 1
	for ; i >= 0; i-- {
		if checker(s[i]) {
			return i
		}
	}

	return i
}
