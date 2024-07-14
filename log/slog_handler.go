package log

import (
	"context"
	"io"
	"log/slog"
)

type SlogHandler struct {
	Handler slog.Handler
	Level   *LogLevel[slog.Level]
}

type SlogHandlerOptions struct {
	// ReplaceAttr see slog.HandlerOptions.ReplaceAttr
	ReplaceAttr func(groups []string, a slog.Attr) slog.Attr
	// Level see slog.HandlerOptions.Level
	Level slog.Leveler
}

// slog.HandlerOptions

func NewSlogHandler(w io.Writer, opts *SlogHandlerOptions) (slog.Handler, LevelToggler) {
	if opts.ReplaceAttr == nil {
		opts.ReplaceAttr = func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.LevelKey {
				a.Value = SlogLevelEncoder(a.Value.Any().(slog.Level))
			}
			return a
		}
	}

	if opts.Level == nil {
		opts.Level = slog.Level(LevelInfo | LevelWarn | LevelError | LevelFatal)
	}

	level := &LogLevel[slog.Level]{opts.Level.Level()}

	return &SlogHandler{slog.NewJSONHandler(w, &slog.HandlerOptions{
			AddSource:   false,
			ReplaceAttr: opts.ReplaceAttr,
			Level:       opts.Level,
		}), level}, func(l Level) {
			level.value = slog.Level(l)
		}
}

func (sh *SlogHandler) Handle(ctx context.Context, r slog.Record) error {
	return sh.Handler.Handle(ctx, r)
}

func (sh *SlogHandler) Enabled(_ context.Context, level slog.Level) bool {
	return (sh.Level.value & level) == level
}

func (sh *SlogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &SlogHandler{sh.Handler.WithAttrs(attrs), sh.Level}
}

func (sh *SlogHandler) WithGroup(name string) slog.Handler {
	return &SlogHandler{sh.Handler.WithGroup(name), sh.Level}
}
