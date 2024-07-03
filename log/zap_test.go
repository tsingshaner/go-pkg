package log

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type levelEnabled[T any] struct {
	fn      func(string, ...T)
	enabled bool
}

type levelTestCases[T any] struct {
	level   Level
	expects []levelEnabled[T]
}

func TestZap(t *testing.T) {
	board := &mockedBoard{}
	core, _ := NewZapCore(
		NewZapJSONEncoder(),
		zapcore.AddSync(board),
		LevelAll,
	)

	logger := NewZapLogger(core)

	assert.Implements(t, (*Logger[zap.Field, zapcore.Level])(nil), logger)

	testCases := []struct {
		logFunc func(string, ...zap.Field)
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

func TestZapLevelToggle(t *testing.T) {
	board := &mockedBoard{}
	core, levelToggler := NewZapCore(
		NewZapJSONEncoder(),
		zapcore.AddSync(board),
		LevelAll,
	)

	logger := NewZapLogger(core)

	testCases := []levelTestCases[zap.Field]{
		{
			LevelTrace | LevelInfo | LevelWarn,
			[]levelEnabled[zap.Field]{
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
			[]levelEnabled[zap.Field]{
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
		levelToggler(tc.level)

		for _, e := range tc.expects {
			board.Flush()
			e.fn("test level")

			assert.Equal(t, e.enabled, board.Size() == 1)
		}
	}
}

func TestZapChild(t *testing.T) {
	board := &mockedBoard{}
	core, _ := NewZapCore(
		NewZapJSONEncoder(),
		zapcore.AddSync(board),
		LevelAll,
	)

	logger := NewZapLogger(core)

	child := logger.Child(zap.Int("pid", 123))

	board.Flush()
	child.Info("test child")

	assert.Equal(t, 1, board.Size())
	assert.Contains(t, string(board.records[0]), "\"msg\":\"test child\"")
	assert.Contains(t, string(board.records[0]), "\"pid\":123")

	board.Flush()
	child2 := child.Child(zap.String("name", "test"))
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

func TestZapGroup(t *testing.T) {
	board := &mockedBoard{}
	core, _ := NewZapCore(
		NewZapJSONEncoder(),
		zapcore.AddSync(board),
		LevelAll,
	)

	logger := NewZapLogger(core)

	grouped := logger.WithGroup("obj")
	grouped.Info("test group", zap.Bool("nested", true))

	assert.Equal(t, 1, board.Size())
	assert.Contains(t, string(board.records[0]), "\"msg\":\"test group\"")
	assert.Contains(t, string(board.records[0]), "\"obj\":{\"nested\":true}")

	board.Flush()

	groupedNested := grouped.WithGroup("obj2")
	groupedNested.Info("test nested group", zap.String("key", "value"))

	assert.Equal(t, 1, board.Size())
	assert.Contains(t, string(board.records[0]), "\"msg\":\"test nested group\"")
	assert.Contains(t, string(board.records[0]), "\"obj\":{\"obj2\":{\"key\":\"value\"}}")

	board.Flush()
	childAndGrouped := grouped.Child(zap.Int("pid", 123)).WithGroup("obj2")
	childAndGrouped.Info("test child and group", zap.Bool("nested", true))
	assert.Equal(t, 1, board.Size())
	assert.Contains(t, string(board.records[0]), "\"msg\":\"test child and group\"")
	assert.Contains(t, string(board.records[0]),
		"\"obj\":{\"pid\":123,\"obj2\":{\"nested\":true}}")
}

func TestZapNamed(t *testing.T) {
	board := &mockedBoard{}
	core, _ := NewZapCore(
		NewZapJSONEncoder(),
		zapcore.AddSync(board),
		LevelAll,
	)

	logger := NewZapLogger(core)

	named := logger.Named("app")
	named.Info("test named", zap.Bool("nested", true))

	assert.Equal(t, 1, board.Size())
	assert.Contains(t, string(board.records[0]), "\"msg\":\"test named\"")
	assert.Contains(t, string(board.records[0]), "\"name\":\"app\"")
	assert.Contains(t, string(board.records[0]), "\"nested\":true")

	board.Flush()

	sub := named.Named("sub")
	sub.Info("test sub named", zap.String("key", "value"))
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

func TestZapSource(t *testing.T) {
	board := &mockedBoard{}
	core, _ := NewZapCore(
		NewZapJSONEncoder(),
		zapcore.AddSync(board),
		LevelAll,
	)

	logger := NewZapLogger(core, zap.AddCaller(), zap.AddCallerSkip(2))

	logger.Info("test source")

	assert.Equal(t, 1, board.Size())
	assert.Contains(t, string(board.records[0]), "\"msg\":\"test source\"")
	assert.Contains(t, string(board.records[0]),
		"\"src\":\"log/zap_test.go:225 github.com/tsingshaner/go-pkg/log.TestZapSource\"")
}

func TestZapStack(t *testing.T) {
	board := &mockedBoard{}
	levelEnabledFn, toggleStackLevel := NewZapLevelFilter(LevelError)
	core, _ := NewZapCore(
		NewZapJSONEncoder(),
		zapcore.AddSync(board),
		LevelAll,
	)
	logger := NewZapLogger(core, zap.AddStacktrace(levelEnabledFn))

	testCases := []levelTestCases[zap.Field]{
		{
			LevelTrace | LevelInfo | LevelWarn,
			[]levelEnabled[zap.Field]{
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
			[]levelEnabled[zap.Field]{
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
					"\"stack\":\"github.com/tsingshaner/go-pkg/log.(*zapLogger).log\\n")
			} else {
				assert.NotContains(t, string(board.records[0]), "\"stack\":")
			}
		}
	}
}

func TestZapSync(t *testing.T) {
	board := &mockedBoard{}
	core, _ := NewZapCore(
		NewZapJSONEncoder(),
		zapcore.AddSync(board),
		LevelAll,
	)
	logger := NewZapLogger(core)

	assert.Nil(t, logger.Sync())
}
