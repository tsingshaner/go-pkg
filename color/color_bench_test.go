package color

import (
	"testing"
)

func BenchmarkRed(b *testing.B) {
	Enable()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Red("hello")
	}
}

func BenchmarkUnsafeMagenta(b *testing.B) {
	Enable()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		UnsafeMagenta("hello")
	}
}

func BenchmarkMulti(b *testing.B) {
	Enable()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Underline(Bold(Red("hello")))
	}
}

func BenchmarkUnsafeMulti(b *testing.B) {
	Enable()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		UnsafeUnderline(UnsafeBold(UnsafeRed("hello")))
	}
}
