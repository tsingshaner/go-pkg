package log

import (
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockedBoard struct {
	records [][]byte
}

func (m *mockedBoard) Write(p []byte) (n int, err error) {
	if string(p) == "@@write error" {
		return 0, errors.New("write error")
	}

	m.records = append(m.records, p)
	return len(p), nil
}

func (m *mockedBoard) Flush() {
	m.records = m.records[:0]
}

func (m *mockedBoard) Size() int {
	return len(m.records)
}

func TestWriter(t *testing.T) {
	board1 := &mockedBoard{}
	board2 := &mockedBoard{}
	writer := NewWriter(board1, board2)

	assert.Implements(t, (*io.Writer)(nil), writer)

	_, _ = writer.Write([]byte("test"))

	assert.Equal(t, 1, board1.Size())
	assert.Equal(t, 1, board2.Size())

	assert.Equal(t, "test", string(board1.records[0]))
	assert.Equal(t, "test", string(board2.records[0]))
}

func TestWriterWithNil(t *testing.T) {
	board := &mockedBoard{}
	writer := NewWriter(board, nil)

	_, _ = writer.Write([]byte("test"))

	assert.Equal(t, 1, board.Size())
	assert.Equal(t, "test", string(board.records[0]))
}

func TestWriterWithError(t *testing.T) {
	board := &mockedBoard{}
	writer := NewWriter(board)

	_, err := writer.Write([]byte("@@write error"))

	assert.Error(t, err)
}
