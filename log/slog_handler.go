package log

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"path/filepath"
)

type SlogHandler struct {
	Handler slog.Handler
	Level   *LogLevel[slog.Level]
}

type SlogHandlerOptions = slog.HandlerOptions

func NewSlogHandler(w io.Writer, opts *SlogHandlerOptions) (slog.Handler, LevelToggler) {
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
					a.Value = SlogLevelEncoder(a.Value.Any().(slog.Level))
				}
				return a
			}
		} else {
			opts.ReplaceAttr = func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == slog.LevelKey {
					a.Value = SlogLevelEncoder(a.Value.Any().(slog.Level))
				}
				return a
			}
		}
	}

	if opts.Level == nil {
		opts.Level = slog.Level(LevelInfo | LevelWarn | LevelError | LevelFatal)
	}

	level := &LogLevel[slog.Level]{opts.Level.Level()}

	return &SlogHandler{slog.NewJSONHandler(w, opts), level}, func(l Level) {
		level.value = slog.Level(l)
	}
}

func (sh *SlogHandler) Handle(ctx context.Context, r slog.Record) error {
	return sh.Handler.Handle(ctx, r)
}

func (sh *SlogHandler) Enabled(_ context.Context, level slog.Level) bool {
	return sh.Level.value&level == level
}

func (sh *SlogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &SlogHandler{sh.Handler.WithAttrs(attrs), sh.Level}
}

func (sh *SlogHandler) WithGroup(name string) slog.Handler {
	return &SlogHandler{sh.Handler.WithGroup(name), sh.Level}
}
