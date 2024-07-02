package log

import (
	"context"
	"io"
	"log/slog"
	"runtime"
	"time"

	"github.com/tsingshaner/go-pkg/util"
)

type Slogger struct {
	logger *slog.Logger
	opts   *Options
	name   string
}

type Options struct {
	addSource bool
	// SkipCaller is the number of stack frames to skip to find the caller.
	SkipCaller int
}

// NewSlog base on go std lib log/slog
func NewSlog(
	w io.Writer,
	slogOpts *SlogHandlerOptions,
	fns ...util.WithFn[Options],
) (Logger[slog.Attr, slog.Level], LevelToggler) {
	loggerOpts := util.BuildWithOpts(&Options{
		addSource:  slogOpts.AddSource,
		SkipCaller: 0,
	}, fns...)

	handler, levelToggler := NewSlogHandler(w, slogOpts)

	return &Slogger{slog.New(handler), loggerOpts, ""}, levelToggler
}

func (s *Slogger) Sync() error { return nil }

func (s *Slogger) Trace(msg string, attrs ...slog.Attr) {
	s.logAttrs(context.Background(), SlogLevelTrace, msg, attrs)
}

func (s *Slogger) Debug(msg string, attrs ...slog.Attr) {
	s.logAttrs(context.Background(), SlogLevelDebug, msg, attrs)
}

func (s *Slogger) Info(msg string, attrs ...slog.Attr) {
	s.logAttrs(context.Background(), SlogLevelInfo, msg, attrs)
}

func (s *Slogger) Warn(msg string, attrs ...slog.Attr) {
	s.logAttrs(context.Background(), SlogLevelWarn, msg, attrs)
}

func (s *Slogger) Error(msg string, attrs ...slog.Attr) {
	s.logAttrs(context.Background(), SlogLevelError, msg, attrs)
}

func (s *Slogger) Fatal(msg string, attrs ...slog.Attr) {
	s.logAttrs(context.Background(), SlogLevelFatal, msg, attrs)
}

func (s *Slogger) Child(attrs ...slog.Attr) Logger[slog.Attr, slog.Level] {
	if len(attrs) == 0 {
		return s
	}

	args := make([]any, 0, len(attrs))
	for _, attr := range attrs {
		args = append(args, attr)
	}

	return &Slogger{s.logger.With(args...), s.opts, s.name}
}

func (s *Slogger) Named(name string) Logger[slog.Attr, slog.Level] {
	if name == "" {
		return s
	}
	if s.name == "" {
		return &Slogger{s.logger, s.opts, name}
	}

	return &Slogger{s.logger, s.opts, s.name + "." + name}
}

func (s *Slogger) Log(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr) {
	s.logAttrs(ctx, level, msg, attrs)
}

func (s *Slogger) WithGroup(name string) Logger[slog.Attr, slog.Level] {
	if name == "" {
		return s
	}

	return &Slogger{s.logger.WithGroup(name), s.opts, s.name}
}

// logAttrs for record callers
func (s *Slogger) logAttrs(ctx context.Context, level slog.Level, msg string, attrs []slog.Attr) {
	if !s.logger.Enabled(context.Background(), level) {
		return
	}
	var pc uintptr
	if s.opts.addSource {
		var pcs [1]uintptr
		// skip [
		//   runtime.Callers,
		//   this function,
		//   this file.Log | LeveledLog,
		//   (this file.Log | LeveledLog)'s caller,
		// ]
		runtime.Callers(s.opts.SkipCaller+3, pcs[:])
		pc = pcs[0]
	}
	r := slog.NewRecord(time.Now(), level, msg, pc)

	if s.name != "" {
		r.AddAttrs(slog.String("name", s.name))
	}

	r.AddAttrs(attrs...)
	if ctx == nil {
		ctx = context.Background()
	}
	_ = s.logger.Handler().Handle(ctx, r)
}
