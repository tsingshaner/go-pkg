package log

import (
	"fmt"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestZapSlog(t *testing.T) {
	board := &mockedBoard{}
	core, _ := NewZapCore(
		NewZapJSONEncoder(),
		zapcore.AddSync(board),
		LevelAll,
	)
	logger := NewZapLog(core, zap.AddCaller(), zap.AddCallerSkip(2))

	assert.Implements(t, (*Logger[slog.Attr, slog.Level])(nil), logger)

	testCases := []struct {
		logFunc func(string, ...slog.Attr)
		level   string
	}{
		{logger.Trace, "trace"},
		{logger.Debug, "debug"},
		{logger.Info, "info"},
		{logger.Warn, "warn"},
		{logger.Error, "error"},
		{logger.Fatal, "fatal"},
	}

	for _, tc := range testCases {
		board.Flush()
		message := fmt.Sprintf("test %s", tc.level)

		tc.logFunc(message)

		assert.Equal(t, 1, board.Size())
		assert.Contains(t, string(board.records[0]), fmt.Sprintf("\"msg\":\"%s\"", message))
		assert.Contains(t, string(board.records[0]), "\"time\":")
		assert.Contains(t, string(board.records[0]), fmt.Sprintf("\"level\":\"%s\"", tc.level))
	}
}

func TestZapSlogLevelToggle(t *testing.T) {
	board := &mockedBoard{}
	core, levelToggler := NewZapCore(
		NewZapJSONEncoder(),
		zapcore.AddSync(board),
		LevelInfo|LevelWarn|LevelError|LevelFatal,
	)
	logger := NewZapLog(core, zap.AddCaller(), zap.AddCallerSkip(2))

	testCases := []levelTestCases{
		{
			LevelTrace | LevelInfo | LevelWarn,
			[]levelEnabled{
				{logger.Trace, true},
				{logger.Debug, false},
				{logger.Info, true},
				{logger.Warn, true},
				{logger.Error, false},
				{logger.Fatal, false},
			},
		},
		{
			LevelDebug | LevelWarn | LevelError,
			[]levelEnabled{
				{logger.Trace, false},
				{logger.Debug, true},
				{logger.Info, false},
				{logger.Warn, true},
				{logger.Error, true},
				{logger.Fatal, false},
			},
		},
		{
			LevelDebug | LevelFatal,
			[]levelEnabled{
				{logger.Trace, false},
				{logger.Debug, true},
				{logger.Info, false},
				{logger.Warn, false},
				{logger.Error, false},
				{logger.Fatal, true},
			},
		},
		{
			LevelSilent,
			[]levelEnabled{
				{logger.Trace, false},
				{logger.Debug, false},
				{logger.Info, false},
				{logger.Warn, false},
				{logger.Error, false},
				{logger.Fatal, false},
			},
		},
	}

	for _, tc := range testCases {
		levelToggler(tc.level)

		for _, e := range tc.expects {
			board.Flush()
			e.fn("test level")

			assert.Equal(t, e.enabled, board.Size() == 1)
		}
	}
}

func TestZapSlogChild(t *testing.T) {
	board := &mockedBoard{}
	core, _ := NewZapCore(
		NewZapJSONEncoder(),
		zapcore.AddSync(board),
		LevelAll,
	)

	logger := NewZapLog(core)

	child := logger.Child(slog.Int("pid", 123))

	board.Flush()
	child.Info("test child")

	assert.Equal(t, 1, board.Size())
	assert.Contains(t, string(board.records[0]), "\"msg\":\"test child\"")
	assert.Contains(t, string(board.records[0]), "\"pid\":123")

	board.Flush()
	child2 := child.Child(slog.String("name", "test"))
	child2.Info("test nesting child")

	assert.Equal(t, 1, board.Size())
	assert.Contains(t, string(board.records[0]), "\"msg\":\"test nesting child\"")
	assert.Contains(t, string(board.records[0]), "\"pid\":123")
	assert.Contains(t, string(board.records[0]), "\"name\":\"test\"")

	board.Flush()
	child3 := logger.Child()
	child3.Info("test empty child")

	assert.Equal(t, 1, board.Size())
	assert.Contains(t, string(board.records[0]), "\"msg\":\"test empty child\"")
	assert.NotContains(t, string(board.records[0]), "\"pid\":")
	assert.NotContains(t, string(board.records[0]), "\"name\":")
}

func TestZapSlogGroup(t *testing.T) {
	board := &mockedBoard{}
	core, _ := NewZapCore(
		NewZapJSONEncoder(),
		zapcore.AddSync(board),
		LevelInfo|LevelWarn|LevelError|LevelFatal,
	)
	logger := NewZapLog(core, zap.AddCaller(), zap.AddCallerSkip(2))

	grouped := logger.WithGroup("obj")
	grouped.Info("test group", slog.Bool("nested", true))

	assert.Equal(t, 1, board.Size())
	assert.Contains(t, string(board.records[0]), "\"msg\":\"test group\"")
	assert.Contains(t, string(board.records[0]), "\"obj\":{\"nested\":true}")

	board.Flush()

	groupedNested := grouped.WithGroup("obj2")
	groupedNested.Info("test nested group", slog.String("key", "value"))

	assert.Equal(t, 1, board.Size())
	assert.Contains(t, string(board.records[0]), "\"msg\":\"test nested group\"")
	assert.Contains(t, string(board.records[0]), "\"obj\":{\"obj2\":{\"key\":\"value\"}}")

	board.Flush()
	childAndGrouped := grouped.Child(slog.Int("pid", 123)).WithGroup("obj2")
	childAndGrouped.Info("test child and group", slog.Bool("nested", true))
	assert.Equal(t, 1, board.Size())
	assert.Contains(t, string(board.records[0]), "\"msg\":\"test child and group\"")
	assert.Contains(t, string(board.records[0]),
		"\"obj\":{\"pid\":123,\"obj2\":{\"nested\":true}}")
}

func TestZapSlogNamed(t *testing.T) {
	board := &mockedBoard{}
	core, _ := NewZapCore(
		NewZapJSONEncoder(),
		zapcore.AddSync(board),
		LevelInfo|LevelWarn|LevelError|LevelFatal,
	)
	logger := NewZapLog(core, zap.AddCaller(), zap.AddCallerSkip(2))

	named := logger.Named("app")
	named.Info("test named", slog.Bool("nested", true))

	assert.Equal(t, 1, board.Size())
	assert.Contains(t, string(board.records[0]), "\"msg\":\"test named\"")
	assert.Contains(t, string(board.records[0]), "\"name\":\"app\"")
	assert.Contains(t, string(board.records[0]), "\"nested\":true")

	board.Flush()

	sub := named.Named("sub")
	sub.Info("test sub named", slog.String("key", "value"))
	assert.Equal(t, 1, board.Size())
	assert.Contains(t, string(board.records[0]), "\"msg\":\"test sub named\"")
	assert.Contains(t, string(board.records[0]), "\"name\":\"app.sub\"")
	assert.Contains(t, string(board.records[0]), "\"key\":\"value\"")

	board.Flush()

	empty := named.Named("")
	empty.Info("test empty named")
	assert.Equal(t, 1, board.Size())
	assert.Contains(t, string(board.records[0]), "\"msg\":\"test empty named\"")
	assert.Contains(t, string(board.records[0]), "\"name\":\"app\"")
}

func TestZapSlogSource(t *testing.T) {
	board := &mockedBoard{}
	core, _ := NewZapCore(
		NewZapJSONEncoder(),
		zapcore.AddSync(board),
		LevelInfo|LevelWarn|LevelError|LevelFatal,
	)
	logger := NewZapLog(core, zap.AddCaller(), zap.AddCallerSkip(2))

	logger.Info("test source")

	assert.Equal(t, 1, board.Size())
	assert.Contains(t, string(board.records[0]), "\"msg\":\"test source\"")
	assert.Contains(t, string(board.records[0]), "\"log/zapslog_test.go:")
	assert.Contains(t, string(board.records[0]),
		"github.com/tsingshaner/go-pkg/log.TestZapSlogSource\"")
}

func TestZapSlogStack(t *testing.T) {
	levelEnabledFn, toggleStackLevel := NewZapLevelFilter(LevelError)
	board := &mockedBoard{}
	core, _ := NewZapCore(
		NewZapJSONEncoder(),
		zapcore.AddSync(board),
		LevelAll,
	)
	logger := NewZapLog(core, zap.AddCaller(), zap.AddCallerSkip(2), zap.AddStacktrace(levelEnabledFn))

	testCases := []levelTestCases{
		{
			LevelTrace | LevelInfo | LevelWarn,
			[]levelEnabled{
				{logger.Trace, true},
				{logger.Debug, false},
				{logger.Info, true},
				{logger.Warn, true},
				{logger.Error, false},
				{logger.Fatal, false},
			},
		},
		{
			LevelDebug | LevelWarn | LevelError,
			[]levelEnabled{
				{logger.Trace, false},
				{logger.Debug, true},
				{logger.Info, false},
				{logger.Warn, true},
				{logger.Error, true},
				{logger.Fatal, false},
			},
		},
	}

	for _, tc := range testCases {
		toggleStackLevel(tc.level)

		for _, e := range tc.expects {
			board.Flush()
			e.fn("test stack")

			assert.Equal(t, 1, board.Size())
			if e.enabled {
				assert.Contains(t, string(board.records[0]),
					"\"stack\":\"github.com/tsingshaner/go-pkg/log.TestZapSlogStack\\n")
			} else {
				assert.NotContains(t, string(board.records[0]), "\"stack\":")
			}
		}
	}
}

func TestZapSlogSync(t *testing.T) {
	board := &mockedBoard{}
	core, _ := NewZapCore(
		NewZapJSONEncoder(),
		zapcore.AddSync(board),
		LevelInfo|LevelWarn|LevelError|LevelFatal,
	)
	logger := NewZapLog(core, zap.AddCaller(), zap.AddCallerSkip(2))

	logger.Info("test sync")
	assert.Nil(t, logger.Sync())
}
