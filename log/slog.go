package log

import (
	"context"
	"io"
	"log/slog"
	"path/filepath"
)

const (
	LevelSilent slog.Level = -1
	LevelTrace  slog.Level = 1 << 0
	LevelDebug  slog.Level = 1 << 1
	LevelInfo   slog.Level = 1 << 2
	LevelWarn   slog.Level = 1 << 3
	LevelError  slog.Level = 1 << 4
	LevelFatal  slog.Level = 1 << 5
)

type SlogHandler struct {
	Handler slog.Handler
	Level   slog.Level
}

func NewSlogHandler(w io.Writer, opts *slog.HandlerOptions) slog.Handler {
	if opts.AddSource && opts.ReplaceAttr == nil {
		opts.ReplaceAttr = func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.SourceKey {
				source := a.Value.Any().(*slog.Source)
				source.File = filepath.Base(source.File)
			}

			if a.Key == slog.LevelKey {
				level := a.Value.Any().(slog.Level)
				switch {
				case level < LevelDebug:
					a.Value = slog.StringValue("trace")
				case level < LevelInfo:
					a.Value = slog.StringValue("debug")
				case level < LevelWarn:
					a.Value = slog.StringValue("info")
				case level < LevelError:
					a.Value = slog.StringValue("warn")
				case level < LevelFatal:
					a.Value = slog.StringValue("error")
				case level < LevelFatal:
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
	println(sh.Level, level)
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
}

func NewSlog(h slog.Handler) *Slogger {
	return &Slogger{logger: slog.New(h)}
}

func (s *Slogger) Trace(msg string, args ...slog.Attr) {
	s.logger.LogAttrs(context.Background(), LevelTrace, msg, args...)
}

func (s *Slogger) Debug(msg string, args ...slog.Attr) {
	s.logger.LogAttrs(context.Background(), LevelDebug, msg, args...)
}

func (s *Slogger) Info(msg string, args ...slog.Attr) {
	s.logger.LogAttrs(context.Background(), LevelInfo, msg, args...)
}

func (s *Slogger) Warn(msg string, args ...slog.Attr) {
	s.logger.LogAttrs(context.Background(), LevelWarn, msg, args...)
}

func (s *Slogger) Error(msg string, args ...slog.Attr) {
	s.logger.LogAttrs(context.Background(), LevelError, msg, args...)
}

func (s *Slogger) Fatal(msg string, args ...slog.Attr) {
	s.logger.LogAttrs(context.Background(), LevelFatal, msg, args...)
}

func (s *Slogger) Child(attrs ...any) *Slogger {
	if len(attrs) == 0 {
		return s
	}

	return &Slogger{logger: s.logger.With(attrs[0])}
}

func (s *Slogger) Log(ctx context.Context, level slog.Level, msg string, args ...slog.Attr) {
	s.logger.LogAttrs(ctx, level, msg, args...)
}

func (s *Slogger) WithGroup(name string) *Slogger {
	if name == "" {
		return s
	}

	return &Slogger{logger: s.logger.WithGroup(name)}
}

func (s *Slogger) LogAttrs(tx context.Context, level slog.Level, msg string, attrs ...slog.Attr) {
	s.logger.LogAttrs(tx, level, msg, attrs...)
}
