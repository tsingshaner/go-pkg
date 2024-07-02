package log

import (
	"context"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewZapJSONEncoder() zapcore.Encoder {
	return zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "name",
		CallerKey:      "src",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stack",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    ZapLevelEncoder,
		EncodeTime:     zapcore.RFC3339NanoTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller: func(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
			src := &strings.Builder{}
			src.WriteString(caller.TrimmedPath())
			src.WriteByte(' ')
			src.WriteString(caller.Function)
			enc.AppendString(src.String())
		},
	})
}

func NewZapLevelFilter(levelBitMask Level) (zap.LevelEnablerFunc, LevelToggler) {
	level := &LogLevel[zapcore.Level]{zapcore.Level(levelBitMask)}

	return zap.LevelEnablerFunc(func(l zapcore.Level) bool { return l&level.value == l }),
		func(l Level) { level.value = zapcore.Level(l) }
}

func NewZapCore(
	enc zapcore.Encoder, ws zapcore.WriteSyncer, l Level,
) (zapcore.Core, LevelToggler) {
	lef, lt := NewZapLevelFilter(l)
	return zapcore.NewCore(enc, ws, lef), lt
}

func NewZapCoreWithFilter(
	enc zapcore.Encoder, ws zapcore.WriteSyncer, filter zap.LevelEnablerFunc,
) zapcore.Core {
	return zapcore.NewCore(enc, ws, filter)
}

type zapLogger struct {
	logger *zap.Logger
}

type LevelEffect Level

func (l LevelEffect) OnWrite(_ *zapcore.CheckedEntry, _ []zapcore.Field) {
	// overwrite default Fatal & Error level will exit the process
}

func NewZapLogger(core zapcore.Core, options ...zap.Option) Logger[zap.Field, zapcore.Level] {
	options = append(options,
		// TODO: apply user custom config
		zap.WithFatalHook(LevelEffect(LevelFatal)),
		zap.WithPanicHook(LevelEffect(LevelFatal)),
	)
	logger := zap.New(core, options...)

	return &zapLogger{logger}
}

func (z *zapLogger) Sync() error { return z.logger.Core().Sync() }

func (z *zapLogger) Trace(msg string, fields ...zap.Field) {
	z.logger.WithOptions(zap.IncreaseLevel(ZapLevelTrace))
	z.log(ZapLevelTrace, msg, fields...)
}

func (z *zapLogger) Debug(msg string, fields ...zap.Field) {
	z.log(ZapLevelDebug, msg, fields...)
}

func (z *zapLogger) Info(msg string, fields ...zap.Field) {
	z.log(ZapLevelInfo, msg, fields...)
}

func (z *zapLogger) Warn(msg string, fields ...zap.Field) {
	z.log(ZapLevelWarn, msg, fields...)
}

func (z *zapLogger) Error(msg string, fields ...zap.Field) {
	z.log(ZapLevelError, msg, fields...)
}

func (z *zapLogger) Fatal(msg string, fields ...zap.Field) {
	z.log(ZapLevelFatal, msg, fields...)
}

func (z *zapLogger) Child(fields ...zap.Field) Logger[zap.Field, zapcore.Level] {
	return &zapLogger{z.logger.With(fields...)}
}

func (z *zapLogger) WithGroup(name string) Logger[zap.Field, zapcore.Level] {
	if name == "" {
		return z
	}

	return &zapLogger{z.logger.With(zap.Namespace(name))}
}

func (z *zapLogger) Named(name string) Logger[zap.Field, zapcore.Level] {
	return &zapLogger{z.logger.Named(name)}
}

func (z *zapLogger) Log(_ context.Context, level zapcore.Level, msg string, fields ...zap.Field) {
	z.log(level, msg, fields...)
}

func (z *zapLogger) log(level zapcore.Level, msg string, fields ...zap.Field) {
	z.logger.Log(level, msg, fields...)
}
