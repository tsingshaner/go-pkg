package log

import (
	"errors"
	"log/slog"
	"os"
	"testing"
	"time"
)

// BenchmarkSlogStruct-8   10000            166139 ns/op             552 B/op         12 allocs/op
func BenchmarkSlogStruct(b *testing.B) {
	logger := NewSlog(os.Stdout, &SlogHandlerOptions{})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("constructed a logger")
		logger.Info("slog constructed a logger",
			slog.String("name", "tsingshaner"),
			slog.Bool("is", true),
			slog.String("err", errors.New("test").Error()),
			slog.Duration("time", time.Nanosecond),
		)

		child := logger.Child(slog.Int("pid", 111))
		child.Info("constructed a logger")
	}
}

// BenchmarkOriginSlog-8 10000            151319 ns/op             480 B/op         11 allocs/op
func BenchmarkOriginSlog(b *testing.B) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("constructed a logger")
		logger.Info("slog constructed a logger",
			slog.String("name", "tsingshaner"),
			slog.Bool("is", true),
			slog.String("err", errors.New("test").Error()),
			slog.Duration("time", time.Nanosecond),
		)

		child := logger.With(slog.Int("pid", 111))
		child.Info("constructed a logger")
	}
}
