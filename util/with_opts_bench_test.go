package util

import "testing"

type options struct {
	A int
	B string
	C bool
}

func withA(a int) WithFn[options] {
	return func(o *options) {
		o.A = a
	}
}

func withB(b string) WithFn[options] {
	return func(o *options) {
		o.B = b
	}
}

func withC(c bool) WithFn[options] {
	return func(o *options) {
		o.C = c
	}
}

// 1000000000               0.2500 ns/op          0 B/op          0 allocs/op
func BenchmarkInline(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		opts := &options{}
		withA(1)(opts)
		withB("hello")(opts)
		withC(true)(opts)
	}
}
func BenchmarkDirect(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		opts := &options{}
		for _, fn := range []WithFn[options]{withA(1), withB("hello"), withC(true)} {
			fn(opts)
		}
	}
}

// 38383417                33.74 ns/op           32 B/op          1 allocs/op
func BenchmarkApplyOpts(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		opts := &options{}
		ApplyOpts(opts, withA(1), withB("hello"), withC(true))
		if opts.A != 1 {
			panic("opts error")
		}
	}
}

// 38024138                33.32 ns/op           32 B/op          1 allocs/op
func BenchmarkBuildWithOpts(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = BuildWithOpts(&options{}, withA(1), withB("hello"), withC(true))
	}
}
