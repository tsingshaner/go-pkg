package log

import (
	"errors"
	"log/slog"
	"testing"
	"time"
)

func BenchmarkCustomSlog(b *testing.B) {
	logger, _ := NewSlog(&mockedBoard{}, &SlogHandlerOptions{})

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

func BenchmarkOriginSlog(b *testing.B) {
	logger := slog.New(slog.NewJSONHandler(&mockedBoard{}, &slog.HandlerOptions{}))
	child := logger.With(slog.Int("pid", 111))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("constructed a logger")
		logger.Info("slog constructed a logger",
			slog.String("name", "tsingshaner"),
			slog.Bool("is", true),
			slog.String("err", errors.New("test").Error()),
			slog.Duration("time", time.Nanosecond),
		)

		child.Info("constructed a logger")
	}
}

func BenchmarkCustomSlogWithSource(b *testing.B) {
	logger, _ := NewSlog(&mockedBoard{}, &SlogHandlerOptions{AddSource: true})
	child := logger.Child(slog.Int("pid", 111))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("constructed a logger")
		logger.Info("slog constructed a logger",
			slog.String("name", "tsingshaner"),
			slog.Bool("is", true),
			slog.String("err", errors.New("test").Error()),
			slog.Duration("time", time.Nanosecond),
		)

		child.Info("constructed a logger")
	}
}

func BenchmarkOriginSlogWithSource(b *testing.B) {
	logger := slog.New(slog.NewJSONHandler(&mockedBoard{}, &slog.HandlerOptions{AddSource: true}))
	child := logger.With(slog.Int("pid", 111))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("constructed a logger")
		logger.Info("slog constructed a logger",
			slog.String("name", "tsingshaner"),
			slog.Bool("is", true),
			slog.String("err", errors.New("test").Error()),
			slog.Duration("time", time.Nanosecond),
		)

		child.Info("constructed a logger")
	}
}

func BenchmarkCustomSlogDisable(b *testing.B) {
	logger, _ := NewSlog(&mockedBoard{}, &SlogHandlerOptions{Level: slog.LevelError})

	b.ResetTimer()
	for i := 0; i < 1000000; i++ {
		logger.Info("constructed a logger")
		logger.Info("slog constructed a logger",
			slog.String("name", "tsingshaner"),
			slog.Bool("is", true),
			slog.String("err", errors.New("test").Error()),
			slog.Duration("time", time.Nanosecond),
		)
	}
}

func BenchmarkOriginSlogDisable(b *testing.B) {
	logger := slog.New(slog.NewJSONHandler(&mockedBoard{}, &slog.HandlerOptions{Level: slog.LevelError}))

	b.ResetTimer()
	for i := 0; i < 1000000; i++ {
		logger.Info("constructed a logger")
		logger.Info("slog constructed a logger",
			slog.String("name", "tsingshaner"),
			slog.Bool("is", true),
			slog.String("err", errors.New("test").Error()),
			slog.Duration("time", time.Nanosecond),
		)
	}
}
