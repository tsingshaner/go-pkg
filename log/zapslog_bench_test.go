package log

import (
	"errors"
	"log/slog"
	"os"
	"testing"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func BenchmarkZapSlog(b *testing.B) {
	core, _ := NewZapCore(
		NewZapJSONEncoder(),
		zapcore.AddSync(&mockedBoard{}),
		LevelInfo|LevelWarn|LevelError|LevelFatal,
	)
	logger := NewZapLog(core)

	defer logger.Sync()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("constructed a logger")
		logger.Info("constructed a logger",
			slog.String("name", "tsingshaner"),
			slog.Bool("is", true),
			slog.String("err", errors.New("test").Error()),
			slog.Duration("time", time.Nanosecond),
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
	logger := zap.New(core)

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

func BenchmarkZapSlogWithSource(b *testing.B) {
	core, _ := NewZapCore(
		NewZapJSONEncoder(),
		zapcore.AddSync(&mockedBoard{}),
		LevelInfo|LevelWarn|LevelError|LevelFatal,
	)
	logger := NewZapLog(core, zap.AddCaller(), zap.AddCallerSkip(2))

	defer logger.Sync()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("constructed a logger")
		logger.Info("constructed a logger",
			slog.String("name", "tsingshaner"),
			slog.Bool("is", true),
			slog.String("err", errors.New("test").Error()),
			slog.Duration("time", time.Nanosecond),
		)
	}
}

func BenchmarkOriginZapWithSource(b *testing.B) {
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

func BenchmarkOriginZapDisable(b *testing.B) {
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(os.Stdout),
		zap.LevelEnablerFunc(func(l zapcore.Level) bool {
			return l >= zapcore.ErrorLevel
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

func BenchmarkZapSlogDisable(b *testing.B) {
	core, _ := NewZapCore(
		NewZapJSONEncoder(),
		zapcore.AddSync(os.Stdout),
		LevelSilent,
	)
	logger := NewZapLog(core, zap.AddCaller(), zap.AddCallerSkip(2))

	defer logger.Sync()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("constructed a logger")
		logger.Error("slog constructed a logger",
			slog.String("name", "tsingshaner"),
			slog.Bool("is", true),
			slog.String("err", errors.New("test").Error()),
			slog.Duration("time", time.Nanosecond),
		)
	}
}
