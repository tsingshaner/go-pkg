package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func TestZapLevelFilter(t *testing.T) {
	levelEnabledFn, toggleLevel := NewZapLevelFilter(LevelError | LevelFatal)

	cases := []struct {
		level   Level
		enabled bool
	}{
		{LevelTrace, false},
		{LevelDebug, false},
		{LevelInfo, false},
		{LevelWarn, false},
		{LevelError, true},
		{LevelFatal, true},
	}

	for _, c := range cases {
		assert.Equal(t, c.enabled, levelEnabledFn(zapcore.Level(c.level)))
	}

	cases = []struct {
		level   Level
		enabled bool
	}{
		{LevelTrace, true},
		{LevelDebug, true},
		{LevelInfo, false},
		{LevelWarn, true},
		{LevelError, false},
		{LevelFatal, true},
	}

	toggleLevel(LevelTrace | LevelDebug | LevelWarn | LevelFatal)
	for _, c := range cases {
		assert.Equal(t, c.enabled, levelEnabledFn(zapcore.Level(c.level)))
	}

	cases = []struct {
		level   Level
		enabled bool
	}{
		{LevelTrace, false},
		{LevelDebug, false},
		{LevelInfo, false},
		{LevelWarn, false},
		{LevelError, false},
		{LevelFatal, false},
	}

	toggleLevel(LevelSilent)
	for _, c := range cases {
		assert.Equal(t, c.enabled, levelEnabledFn(zapcore.Level(c.level)))
	}
}
