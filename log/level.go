package log

import (
	"log/slog"
	"strconv"

	"go.uber.org/zap/zapcore"
	"golang.org/x/exp/constraints"
)

type Level = uint8
type LevelToggler func(Level)

// LogLevel is a wrapper for dynamically changing log level.
type LogLevel[T constraints.Signed] struct {
	value T
}

const (
	LevelSilent Level = 0
	LevelTrace  Level = 1 << 0
	LevelDebug  Level = 1 << 1
	LevelInfo   Level = 1 << 2
	LevelWarn   Level = 1 << 3
	LevelError  Level = 1 << 4
	LevelFatal  Level = 1 << 5
	LevelAll    Level = LevelTrace | LevelDebug | LevelInfo | LevelWarn | LevelError | LevelFatal
)

const (
	Trace = "trace"
	Debug = "debug"
	Info  = "info"
	Warn  = "warn"
	Error = "error"
	Fatal = "fatal"
	Panic = "panic"
)

const (
	SlogLevelSilent = slog.Level(LevelSilent)
	SlogLevelTrace  = slog.Level(LevelTrace)
	SlogLevelDebug  = slog.Level(LevelDebug)
	SlogLevelInfo   = slog.Level(LevelInfo)
	SlogLevelWarn   = slog.Level(LevelWarn)
	SlogLevelError  = slog.Level(LevelError)
	SlogLevelFatal  = slog.Level(LevelFatal)
	SlogLevelAll    = slog.Level(LevelAll)
)

var (
	slogLevelStringTrace = slog.StringValue(Trace)
	slogLevelStringDebug = slog.StringValue(Debug)
	slogLevelStringInfo  = slog.StringValue(Info)
	slogLevelStringWarn  = slog.StringValue(Warn)
	slogLevelStringError = slog.StringValue(Error)
	slogLevelStringFatal = slog.StringValue(Fatal)
)

func SlogLevelEncoder(level slog.Level) slog.Value {
	switch level {
	case SlogLevelTrace:
		return slogLevelStringTrace
	case SlogLevelDebug:
		return slogLevelStringDebug
	case SlogLevelInfo:
		return slogLevelStringInfo
	case SlogLevelWarn:
		return slogLevelStringWarn
	case SlogLevelError:
		return slogLevelStringError
	case SlogLevelFatal:
		return slogLevelStringFatal
	default:
		return slog.StringValue("level(" + strconv.Itoa(int(level)) + ")")
	}
}

const (
	ZapLevelSilent = zapcore.Level(LevelSilent)
	ZapLevelTrace  = zapcore.Level(LevelTrace)
	ZapLevelDebug  = zapcore.Level(LevelDebug)
	ZapLevelInfo   = zapcore.Level(LevelInfo)
	ZapLevelWarn   = zapcore.Level(LevelWarn)
	ZapLevelError  = zapcore.Level(LevelError)
	ZapLevelFatal  = zapcore.Level(LevelFatal)
	ZapLevelAll    = zapcore.Level(LevelAll)
)

func ZapLevelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	switch l {
	case ZapLevelTrace:
		enc.AppendString(Trace)
	case ZapLevelDebug:
		enc.AppendString(Debug)
	case ZapLevelInfo:
		enc.AppendString(Info)
	case ZapLevelWarn:
		enc.AppendString(Warn)
	case ZapLevelError:
		enc.AppendString(Error)
	case ZapLevelFatal:
		enc.AppendString(Fatal)
	default:
		enc.AppendString("level(" + strconv.Itoa(int(l)) + ")")
	}
}
