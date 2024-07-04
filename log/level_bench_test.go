package log

import (
	"log/slog"
	"testing"
)

func BenchmarkSlogLevelEncoder(b *testing.B) {
	levels := []slog.Level{
		SlogLevelTrace,
		SlogLevelDebug,
		SlogLevelInfo,
		SlogLevelWarn,
		SlogLevelError,
		SlogLevelFatal,
		SlogLevelAll,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, level := range levels {
			SlogLevelEncoder(level)
		}
	}
}
