package log

import (
	"errors"
	"testing"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func BenchmarkZapLogger(b *testing.B) {
	core, _ := NewZapCore(
		NewZapJSONEncoder(),
		zapcore.AddSync(&mockedBoard{}),
		LevelInfo|LevelWarn|LevelError|LevelFatal,
	)
	logger := NewZapLogger(core, zap.AddCaller(), zap.AddCallerSkip(2))

	defer logger.Sync()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("constructed a logger")
		logger.Info("constructed a logger",
			zap.String("name", "tsingshaner"),
			zap.Bool("is", true),
			zap.Error(errors.New("test")),
			zap.Duration("time", time.Nanosecond),
		)
	}
}

func BenchmarkOriginZap(b *testing.B) {
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(&mockedBoard{}),
		zap.LevelEnablerFunc(func(l zapcore.Level) bool {
			return l >= zapcore.InfoLevel
		}),
	)
	logger := zap.New(core, zap.AddCaller())
	defer logger.Sync()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("constructed a logger")
		logger.Info("constructed a logger",
			zap.String("name", "tsingshaner"),
			zap.Bool("is", true),
			zap.Error(errors.New("test")),
			zap.Duration("time", time.Nanosecond),
		)
	}
}
