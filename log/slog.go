package log

import (
	"context"
	"io"
	"log/slog"
	"path/filepath"
	"runtime"
	"time"
)

const (
	LevelSilent slog.Level = -1
	LevelTrace  slog.Level = 1 << 0
	LevelDebug  slog.Level = 1 << 1
	LevelInfo   slog.Level = 1 << 2
	LevelWarn   slog.Level = 1 << 3
	LevelError  slog.Level = 1 << 4
	LevelFatal  slog.Level = 1 << 5
	LevelAll    slog.Level = LevelTrace | LevelDebug | LevelInfo | LevelWarn | LevelError | LevelFatal
)

type SlogHandler struct {
	Handler slog.Handler
	Level   slog.Level
}

type SlogHandlerOptions = slog.HandlerOptions

func NewSlogHandler(w io.Writer, opts *SlogHandlerOptions) slog.Handler {
	if opts.AddSource && opts.ReplaceAttr == nil {
		opts.ReplaceAttr = func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.SourceKey {
				source := a.Value.Any().(*slog.Source)
				source.File = filepath.Base(source.File)
			}

			if a.Key == slog.LevelKey {
				level := a.Value.Any().(slog.Level)
				switch level {
				case LevelTrace:
					a.Value = slog.StringValue("trace")
				case LevelDebug:
					a.Value = slog.StringValue("debug")
				case LevelInfo:
					a.Value = slog.StringValue("info")
				case LevelWarn:
					a.Value = slog.StringValue("warn")
				case LevelError:
					a.Value = slog.StringValue("error")
				case LevelFatal:
					a.Value = slog.StringValue("fatal")
				}
			}
			return a
		}
	}

	if opts.Level == nil {
		opts.Level = LevelInfo | LevelWarn | LevelError | LevelFatal
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
}

type Options struct {
	addSource bool
	// SkipCaller is the number of stack frames to skip to find the caller.
	SkipCaller int
}

// NewSlog base on go std lib log/slog
func NewSlog(w io.Writer, slogOpts *SlogHandlerOptions, fns ...func(*Options)) *Slogger {
	loggerOpts := &Options{
		addSource:  slogOpts.AddSource,
		SkipCaller: 0,
	}
	for _, fn := range fns {
		fn(loggerOpts)
	}

	return &Slogger{slog.New(NewSlogHandler(w, slogOpts)), loggerOpts}
}

func (s *Slogger) Trace(msg string, attrs ...slog.Attr) {
	s.logAttrs(context.Background(), LevelTrace, msg, attrs)
}

func (s *Slogger) Debug(msg string, attrs ...slog.Attr) {
	s.logAttrs(context.Background(), LevelDebug, msg, attrs)
}

func (s *Slogger) Info(msg string, attrs ...slog.Attr) {
	s.logAttrs(context.Background(), LevelInfo, msg, attrs)
}

func (s *Slogger) Warn(msg string, attrs ...slog.Attr) {
	s.logAttrs(context.Background(), LevelWarn, msg, attrs)
}

func (s *Slogger) Error(msg string, attrs ...slog.Attr) {
	s.logAttrs(context.Background(), LevelError, msg, attrs)
}

func (s *Slogger) Fatal(msg string, attrs ...slog.Attr) {
	s.logAttrs(context.Background(), LevelFatal, msg, attrs)
}

func (s *Slogger) Child(attrs ...slog.Attr) *Slogger {
	if len(attrs) == 0 {
		return s
	}

	args := make([]any, 0, len(attrs))
	for _, attr := range attrs {
		args = append(args, attr)
	}

	return &Slogger{s.logger.With(args...), s.opts}
}

func (s *Slogger) Log(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr) {
	s.logAttrs(ctx, level, msg, attrs)
}

func (s *Slogger) WithGroup(name string) *Slogger {
	if name == "" {
		return s
	}

	return &Slogger{s.logger.WithGroup(name), s.opts}
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
	r.AddAttrs(attrs...)
	if ctx == nil {
		ctx = context.Background()
	}
	_ = s.logger.Handler().Handle(ctx, r)
}
