package log

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"path/filepath"
	"runtime"
	"time"

	"github.com/tsingshaner/go-pkg/util"
)

type SlogHandler struct {
	Handler slog.Handler
	Level   slog.Level
}

type SlogHandlerOptions = slog.HandlerOptions

func SlogLevelTag(level slog.Level) slog.Value {
	switch level {
	case SlogLevelTrace:
		return slog.StringValue("trace")
	case SlogLevelDebug:
		return slog.StringValue("debug")
	case SlogLevelInfo:
		return slog.StringValue("info")
	case SlogLevelWarn:
		return slog.StringValue("warn")
	case SlogLevelError:
		return slog.StringValue("error")
	case SlogLevelFatal:
		return slog.StringValue("fatal")
	}

	return slog.StringValue(level.String())
}

func NewSlogHandler(w io.Writer, opts *SlogHandlerOptions) slog.Handler {
	if opts.ReplaceAttr == nil {
		if opts.AddSource {
			opts.ReplaceAttr = func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == slog.SourceKey {
					a.Key = "src"
					source := a.Value.Any().(*slog.Source)
					source.File = filepath.Base(source.File)
					a.Value = slog.StringValue(fmt.Sprintf("%s:%d %s", source.File, source.Line, source.Function))
				}

				if a.Key == slog.LevelKey {
					a.Value = SlogLevelTag(a.Value.Any().(slog.Level))
				}
				return a
			}
		} else {
			opts.ReplaceAttr = func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == slog.LevelKey {
					a.Value = SlogLevelTag(a.Value.Any().(slog.Level))
				}
				return a
			}
		}
	}

	if opts.Level == nil {
		opts.Level = slog.Level(LevelInfo | LevelWarn | LevelError | LevelFatal)
	}

	return &SlogHandler{
		Handler: slog.NewJSONHandler(w, opts),
		Level:   opts.Level.Level(),
	}
}

func (sh *SlogHandler) Handle(ctx context.Context, r slog.Record) error {
	return sh.Handler.Handle(ctx, r)
}

func (sh *SlogHandler) Enabled(_ context.Context, level slog.Level) bool {
	return sh.Level&level == level
}

func (sh *SlogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &SlogHandler{sh.Handler.WithAttrs(attrs), sh.Level}
}

func (sh *SlogHandler) WithGroup(name string) slog.Handler {
	return &SlogHandler{sh.Handler.WithGroup(name), sh.Level}
}

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
) Logger[slog.Attr, slog.Level] {
	loggerOpts := util.BuildWithOpts(&Options{
		addSource:  slogOpts.AddSource,
		SkipCaller: 0,
	}, fns...)

	return &Slogger{slog.New(NewSlogHandler(w, slogOpts)), loggerOpts, ""}
}

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
