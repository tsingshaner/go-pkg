package log

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"path/filepath"
	"time"

	"github.com/tsingshaner/go-pkg/util"
	"github.com/tsingshaner/go-pkg/util/runtime"
)

type slogger struct {
	logger *slog.Logger
	opts   *Options
	name   string
}

type Options struct {
	// AddSource is whether to add source code information to the log, key is 'src'.
	AddSource bool
	// SkipCaller is the number of stack frames to skip to find the caller.
	SkipCaller int
	// StackTrace is the level of enable stack trace to log.
	StackTrace slog.Level
}

// NewSlog base on go std lib log/slog
func NewSlog(
	w io.Writer,
	slogOpts *SlogHandlerOptions,
	fns ...util.WithFn[Options],
) (Logger[slog.Attr, slog.Level], LevelToggler) {
	loggerOpts := util.BuildWithOpts(&Options{false, 0, slog.Level(LevelSilent)}, fns...)

	handler, levelToggler := NewSlogHandler(w, slogOpts)

	return &slogger{slog.New(handler), loggerOpts, ""}, levelToggler
}

func (s *slogger) Sync() error { return nil }

func (s *slogger) Trace(msg string, attrs ...slog.Attr) {
	s.logAttrs(context.Background(), SlogLevelTrace, msg, attrs)
}

func (s *slogger) Debug(msg string, attrs ...slog.Attr) {
	s.logAttrs(context.Background(), SlogLevelDebug, msg, attrs)
}

func (s *slogger) Info(msg string, attrs ...slog.Attr) {
	s.logAttrs(context.Background(), SlogLevelInfo, msg, attrs)
}

func (s *slogger) Warn(msg string, attrs ...slog.Attr) {
	s.logAttrs(context.Background(), SlogLevelWarn, msg, attrs)
}

func (s *slogger) Error(msg string, attrs ...slog.Attr) {
	s.logAttrs(context.Background(), SlogLevelError, msg, attrs)
}

func (s *slogger) Fatal(msg string, attrs ...slog.Attr) {
	s.logAttrs(context.Background(), SlogLevelFatal, msg, attrs)
}

func (s *slogger) Child(attrs ...slog.Attr) Slog {
	if len(attrs) == 0 {
		return s
	}

	args := make([]any, 0, len(attrs))
	for _, attr := range attrs {
		args = append(args, attr)
	}

	return &slogger{s.logger.With(args...), s.opts, s.name}
}

func (s *slogger) WithOptions(opts *ChildLoggerOptions) Slog {
	return &slogger{
		s.logger,
		&Options{opts.AddSource, s.opts.SkipCaller + opts.SkipCaller, slog.Level(opts.StackTrace)},
		s.name,
	}
}

func (s *slogger) Named(name string) Slog {
	if name == "" {
		return s
	}
	if s.name == "" {
		return &slogger{s.logger, s.opts, name}
	}

	return &slogger{s.logger, s.opts, s.name + "." + name}
}

func (s *slogger) Log(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr) {
	s.logAttrs(ctx, level, msg, attrs)
}

func (s *slogger) WithGroup(name string) Slog {
	if name == "" {
		return s
	}

	return &slogger{s.logger.WithGroup(name), s.opts, s.name}
}

func (s *slogger) Enabled(level slog.Level) bool {
	return s.logger.Enabled(context.Background(), level)
}

// logAttrs for record callers
func (s *slogger) logAttrs(ctx context.Context, level slog.Level, msg string, attrs []slog.Attr) {
	if !s.logger.Enabled(context.Background(), level) {
		return
	}
	var pc uintptr
	r := slog.NewRecord(time.Now(), level, msg, pc)

	addStack := (s.opts.StackTrace & level) == level
	if s.opts.AddSource || addStack {
		caller, stack, _ := runtime.GetStackTrace(s.opts.AddSource, addStack, s.opts.SkipCaller+2)

		if s.opts.AddSource {
			r.AddAttrs(slog.String("src", fmt.Sprintf("%s:%d %s",
				filepath.Base(caller.File), caller.Line, caller.Function)))
		}

		if addStack {
			r.AddAttrs(slog.String("stack", stack))
		}
	}

	if s.name != "" {
		r.AddAttrs(slog.String("name", s.name))
	}

	r.AddAttrs(attrs...)
	if ctx == nil {
		ctx = context.Background()
	}
	_ = s.logger.Handler().Handle(ctx, r)
}
