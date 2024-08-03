package util

func Noop(...any) {}

func Pick[T any](conditional bool, a, b T) T {
	if conditional {
		return a
	}
	return b
}
