package log

import "context"

type LoggerFeature[Attr any, Level any] interface {
	// Child
	// - [zh] 创建一个新的记录器，均附带指定的属性。
	// - [en] creates a new logger, all with the specified attributes.
	Child(...Attr) Logger[Attr, Level]
	// Named { "name": "name1.name2.name3" }
	// - [zh] 创建一个新的具有给定名称的记录器，连接到父级后。
	// - [en] creates a new logger with the given name, connected to the parent after.
	Named(name string) Logger[Attr, Level]
	// WithGroup { "level": "info", "name": { "key": "value" } }
	// - [zh] 创建子 logger, 之后的属性添加到 name 字段下。
	// - [en] create a child logger, the following attributes are added to the name field.
	WithGroup(name string) Logger[Attr, Level]
	// Sync only for zap if you used
	// see [zap.Logger.Sync](https://pkg.go.dev/go.uber.org/zap#Logger.Sync)
	Sync() error
}

type Logger[Attr any, Level any] interface {
	LoggerFeature[Attr, Level]

	Log(ctx context.Context, level Level, msg string, attrs ...Attr)
	Trace(msg string, args ...Attr)
	Debug(msg string, args ...Attr)
	Info(msg string, args ...Attr)
	Warn(msg string, args ...Attr)
	Error(msg string, args ...Attr)
	Fatal(msg string, args ...Attr)
}
