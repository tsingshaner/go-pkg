package log

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func TestLevelEncoder(t *testing.T) {
	testCases := []struct {
		level        Level
		expectedName string
	}{
		{LevelTrace, "trace"},
		{LevelDebug, "debug"},
		{LevelInfo, "info"},
		{LevelWarn, "warn"},
		{LevelError, "error"},
		{LevelFatal, "fatal"},
		{3, "level(3)"},
	}

	for _, tc := range testCases {
		t.Run("slog/"+tc.expectedName, func(t *testing.T) {
			assert.Equal(t, tc.expectedName, SlogLevelEncoder(slog.Level(tc.level)).String())
		})
	}

	mockedEnc := new(mockedZapPrimitiveArrayEncoder)

	for _, tc := range testCases {
		t.Run("zap/"+tc.expectedName, func(t *testing.T) {
			mockedEnc.record = mockedEnc.record[:0]
			ZapLevelEncoder(zapcore.Level(tc.level), mockedEnc)
			assert.Equal(t, tc.expectedName, mockedEnc.record[0])
		})

	}
}

type mockedZapPrimitiveArrayEncoder struct {
	zapcore.PrimitiveArrayEncoder
	record []string
}

func (m *mockedZapPrimitiveArrayEncoder) AppendString(s string) {
	m.record = append(m.record, s)
}
