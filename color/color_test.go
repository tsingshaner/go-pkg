package color

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const text = "hello"

func TestRed(t *testing.T) {
	redString := Red(text)

	assert.Contains(t, redString, text)
}

func TestMulti(t *testing.T) {
	multiString := Underline(Bold(Red(text)))

	assert.Contains(t, multiString, text)
}
