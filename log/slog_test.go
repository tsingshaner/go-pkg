package log

import (
	"fmt"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

type levelEnabled struct {
	fn      func(string, ...slog.Attr)
	enabled bool
}

type levelTestCases struct {
	level   Level
	expects []levelEnabled
}

func TestSlog(t *testing.T) {
	board := &mockedBoard{}
	logger, _ := NewSlog(board, &SlogHandlerOptions{
		Level: SlogLevelAll,
	})

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
		t.Run(tc.level, func(t *testing.T) {
			board.Flush()
			message := fmt.Sprintf("test %s", tc.level)

			tc.logFunc(message)

			assert.Equal(t, 1, board.Size())
			assert.Contains(t, string(board.records[0]), fmt.Sprintf("\"msg\":\"%s\"", message))
			assert.Contains(t, string(board.records[0]), "\"time\":")
			assert.Contains(t, string(board.records[0]), fmt.Sprintf("\"level\":\"%s\"", tc.level))
		})
	}
}

func TestSlogLevelToggle(t *testing.T) {
	board := &mockedBoard{}
	logger, levelToggler := NewSlog(board, &SlogHandlerOptions{
		Level: SlogLevelAll,
	})

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
		t.Run(fmt.Sprintf("level(%d)", tc.level), func(t *testing.T) {
			levelToggler(tc.level)

			for _, e := range tc.expects {
				board.Flush()
				e.fn("test level")

				assert.Equal(t, e.enabled, board.Size() == 1)
			}
		})

	}
}

func TestSlogChild(t *testing.T) {
	board := &mockedBoard{}
	logger, _ := NewSlog(board, &SlogHandlerOptions{
		Level: SlogLevelAll,
	})

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

func TestSlogGroup(t *testing.T) {
	board := &mockedBoard{}
	logger, _ := NewSlog(board, &SlogHandlerOptions{
		Level: SlogLevelAll,
	})

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

func TestSlogNamed(t *testing.T) {
	board := &mockedBoard{}
	logger, _ := NewSlog(board, &SlogHandlerOptions{
		Level: SlogLevelAll,
	})

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

func TestSlogSource(t *testing.T) {
	board := &mockedBoard{}
	logger, _ := NewSlog(board, &SlogHandlerOptions{}, func(o *Options) {
		o.AddSource = true
	})

	logger.Info("test source")

	assert.Equal(t, 1, board.Size())
	assert.Contains(t, string(board.records[0]), "\"msg\":\"test source\"")
	assert.Contains(t, string(board.records[0]), "\"src\":\"slog_test.go:")
	assert.Contains(t, string(board.records[0]),
		" github.com/tsingshaner/go-pkg/log.TestSlogSource\"")
}

func TestSlogSync(t *testing.T) {
	board := &mockedBoard{}
	logger, _ := NewSlog(board, &SlogHandlerOptions{})

	assert.Nil(t, logger.Sync())
}

func TestSlogChildWithOptions(t *testing.T) {
	board := &mockedBoard{}
	logger, _ := NewSlog(board, &SlogHandlerOptions{}, func(o *Options) {
		o.AddSource = true
	})

	child := logger.WithOptions(&ChildLoggerOptions{
		AddSource:  false,
		SkipCaller: 0,
		StackTrace: LevelError,
	})

	logger.Info("parent info")
	assert.Equal(t, 1, board.Size())
	assert.Contains(t, string(board.records[0]), "\"src\":\"")

	board.Flush()
	child.Info("child info")

	assert.Equal(t, 1, board.Size())
	assert.NotContains(t, string(board.records[0]), "\"src\":\"")

	child = logger.WithOptions(&ChildLoggerOptions{
		AddSource:  false,
		SkipCaller: 0,
		StackTrace: LevelInfo,
	})

	board.Flush()
	logger.Info("parent info without stack")

	assert.Equal(t, 1, board.Size())
	assert.NotContains(t, string(board.records[0]), "\"stack\":\"")

	board.Flush()
	child.Info("child info with stack")

	assert.Equal(t, 1, board.Size())
	assert.Contains(t, string(board.records[0]), "\"stack\":\"")
	assert.NotContains(t, string(board.records[0]), "\"src\":\"")
}
