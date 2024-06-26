package log

import "context"

type Logger[Attr any, Level any] interface {
	Trace(msg string, args ...Attr)
	Debug(msg string, args ...Attr)
	Info(msg string, args ...Attr)
	Warn(msg string, args ...Attr)
	Error(msg string, args ...Attr)
	Fatal(msg string, args ...Attr)
	Child(...Attr) Logger[Attr, Level]
	WithGroup(name string) Logger[Attr, Level]
	// Named { "name": "name1.name2.name3" }
	// - [zh] 创建一个新的具有给定名称的记录器，连接到父级后。
	// - [en] creates a new logger with the given name, connected to the parent after.
	Named(name string) Logger[Attr, Level]
	Log(ctx context.Context, level Level, msg string, attrs ...Attr)
}
