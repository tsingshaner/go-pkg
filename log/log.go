package log

import (
	"context"
	"log/slog"

	"go.uber.org/zap"
	"golang.org/x/exp/constraints"
)

type Attr interface {
	slog.Attr | zap.Field
}

type LoggerFeature[T Attr, Level constraints.Signed] interface {
	// Child
	// - [zh] 创建一个新的记录器，均附带指定的属性。
	// - [en] creates a new logger, all with the specified attributes.
	Child(...T) Logger[T, Level]
	Enabled(Level) bool
	// Named { "name": "name1.name2.name3" }
	// - [zh] 创建一个新的具有给定名称的记录器，连接到父级后。
	// - [en] creates a new logger with the given name, connected to the parent after.
	Named(name string) Logger[T, Level]
	// WithGroup { "level": "info", "name": { "key": "value" } }
	// - [zh] 创建子 logger, 之后的属性添加到 name 字段下。
	// - [en] create a child logger, the following attributes are added to the name field.
	WithGroup(name string) Logger[T, Level]
	// Sync only for zap if you used
	// see [zap.Logger.Sync](https://pkg.go.dev/go.uber.org/zap#Logger.Sync)
	Sync() error
}

type Logger[T Attr, Level constraints.Signed] interface {
	LoggerFeature[T, Level]

	Log(ctx context.Context, level Level, msg string, attrs ...T)
	Trace(msg string, args ...T)
	Debug(msg string, args ...T)
	Info(msg string, args ...T)
	Warn(msg string, args ...T)
	Error(msg string, args ...T)
	Fatal(msg string, args ...T)
}

type Slog = Logger[slog.Attr, slog.Level]
