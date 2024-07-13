package log

import (
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

	return zap.LevelEnablerFunc(func(l zapcore.Level) bool { return (l & level.value) == l }),
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

type LevelEffect Level

func (l LevelEffect) OnWrite(_ *zapcore.CheckedEntry, _ []zapcore.Field) {
	// overwrite default Fatal & Error level will exit the process
}
