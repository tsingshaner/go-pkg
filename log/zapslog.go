package log

import (
	"context"
	"log/slog"

	"github.com/tsingshaner/go-pkg/util/slices"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapLog struct {
	logger *zap.Logger
}

func NewZapLog(core zapcore.Core, options ...zap.Option) Slog {
	options = append(options,
		// TODO: apply user custom config
		zap.WithFatalHook(LevelEffect(LevelFatal)),
		zap.WithPanicHook(LevelEffect(LevelFatal)),
	)
	logger := zap.New(core, options...)

	return &zapLog{logger}
}

func (z *zapLog) Sync() error { return z.logger.Core().Sync() }

func (z *zapLog) Trace(msg string, fields ...slog.Attr) {
	z.logger.WithOptions(zap.IncreaseLevel(ZapLevelTrace))
	z.log(ZapLevelTrace, msg, fields...)
}

func (z *zapLog) Debug(msg string, fields ...slog.Attr) {
	z.log(ZapLevelDebug, msg, fields...)
}

func (z *zapLog) Info(msg string, fields ...slog.Attr) {
	z.log(ZapLevelInfo, msg, fields...)
}

func (z *zapLog) Warn(msg string, fields ...slog.Attr) {
	z.log(ZapLevelWarn, msg, fields...)
}

func (z *zapLog) Error(msg string, fields ...slog.Attr) {
	z.log(ZapLevelError, msg, fields...)
}

func (z *zapLog) Fatal(msg string, fields ...slog.Attr) {
	z.log(ZapLevelFatal, msg, fields...)
}

func (z *zapLog) Child(fields ...slog.Attr) Slog {
	return &zapLog{z.logger.With(convertSlogToZapFields(fields)...)}
}

func (z *zapLog) WithGroup(name string) Slog {
	if name == "" {
		return z
	}

	return &zapLog{z.logger.With(zap.Namespace(name))}
}

func (z *zapLog) Named(name string) Slog {
	return &zapLog{z.logger.Named(name)}
}

func (z *zapLog) Enabled(level slog.Level) bool {
	return z.logger.Core().Enabled(zapcore.Level(level))
}

func (z *zapLog) Log(_ context.Context, level slog.Level, msg string, fields ...slog.Attr) {
	z.log(zapcore.Level(level), msg, fields...)
}

func (z *zapLog) log(level zapcore.Level, msg string, fields ...slog.Attr) {
	if z.logger.Core().Enabled(level) {
		z.logger.Log(level, msg, convertSlogToZapFields(fields)...)
	}
}

func convertSlogToZapFields(fields []slog.Attr) []zap.Field {
	return slices.Map(fields, convertAttrToField)
}

func convertAttrToField(attr slog.Attr) zapcore.Field {
	if attr.Equal(slog.Attr{}) {
		// Ignore empty attrs.
		return zap.Skip()
	}

	switch attr.Value.Kind() {
	case slog.KindBool:
		return zap.Bool(attr.Key, attr.Value.Bool())
	case slog.KindDuration:
		return zap.Duration(attr.Key, attr.Value.Duration())
	case slog.KindFloat64:
		return zap.Float64(attr.Key, attr.Value.Float64())
	case slog.KindInt64:
		return zap.Int64(attr.Key, attr.Value.Int64())
	case slog.KindString:
		return zap.String(attr.Key, attr.Value.String())
	case slog.KindTime:
		return zap.Time(attr.Key, attr.Value.Time())
	case slog.KindUint64:
		return zap.Uint64(attr.Key, attr.Value.Uint64())
	case slog.KindGroup:
		if attr.Key == "" {
			// Inlines recursively.
			return zap.Inline(groupObject(attr.Value.Group()))
		}
		return zap.Object(attr.Key, groupObject(attr.Value.Group()))
	case slog.KindLogValuer:
		return convertAttrToField(slog.Attr{
			Key:   attr.Key,
			Value: attr.Value.Resolve(),
		})
	default:
		return zap.Any(attr.Key, attr.Value.Any())
	}
}

type groupObject []slog.Attr

func (gs groupObject) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	for _, attr := range gs {
		convertAttrToField(attr).AddTo(enc)
	}
	return nil
}
