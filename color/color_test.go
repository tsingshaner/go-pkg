package color

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const text = "hello"

func TestEnable(t *testing.T) {
	defaultValue := IsEnabled()

	Enable()
	assert.True(t, IsEnabled())

	Disable()
	assert.False(t, IsEnabled())

	ResetEnabled()
	assert.Equal(t, defaultValue, IsEnabled())

	Enable()
	ResetEnabled()
	assert.Equal(t, defaultValue, IsEnabled())
}

func TestSingle(t *testing.T) {
	redString := Red(text)

	assert.Contains(t, redString, text)
	assert.Contains(t, redString, "[31m")

	boldString := Bold(text)

	assert.Contains(t, boldString, text)
	assert.Contains(t, boldString, "[1m")
}

func TestMulti(t *testing.T) {
	multiString := Underline(Bold(Red(text)))

	assert.Contains(t, multiString, text)
}

func TestUnsafe(t *testing.T) {
	redString := UnsafeRed(text)

	assert.Contains(t, redString, text)
	assert.Contains(t, redString, "[31m")

	boldString := UnsafeBold(text)

	assert.Contains(t, boldString, text)
	assert.Contains(t, boldString, "[1m")
}

func TestUnsafeMulti(t *testing.T) {
	multiString := UnsafeUnderline(UnsafeBold(UnsafeRed(text)))

	assert.Contains(t, multiString, text)
	assert.Contains(t, multiString, "[31m")
	assert.Contains(t, multiString, "[1m")
	assert.Contains(t, multiString, "[4m")
}

func TestUnsafeWithFail(t *testing.T) {
	mixinString := UnsafeRed(" red " + Green(" green ") + " no color ")

	assert.Equal(t, "\x1b[31m red \x1b[32m green \x1b[39m no color \x1b[39m", mixinString)
}

func TestSafeFormat(t *testing.T) {
	mixinString := Red(" red " + Green(" green ") + " red ")

	assert.Equal(t, "\x1b[31m red \x1b[32m green \x1b[31m red \x1b[39m", mixinString)
}

func TestBoldAndDimIsSafe(t *testing.T) {
	mixinString := Bold(" bold " + Dim(" dim ") + " bold ")
	assert.Equal(t, "\x1b[1m bold \x1b[2m dim \x1b[22m\x1b[1m bold \x1b[22m", mixinString)

	mixinString = Dim(" dim " + Bold(" bold ") + " dim ")
	assert.Equal(t, "\x1b[2m dim \x1b[1m bold \x1b[22m\x1b[2m dim \x1b[22m", mixinString)
}

func TestMultiStyle(t *testing.T) {
	mixinString := Underline(Bold(Dim(" dim " + Red(" red ") + " dim ")))
	assert.Equal(t,
		"\x1b[4m\x1b[1m\x1b[2m dim \x1b[31m red \x1b[39m dim \x1b[22m\x1b[1m\x1b[22m\x1b[24m",
		mixinString)

	mixinString = Dim(Red(" red " + Bold(" bold ") + Yellow(" yellow ")))
	assert.Equal(t,
		"\x1b[2m\x1b[31m red \x1b[1m bold \x1b[22m\x1b[2m\x1b[33m yellow \x1b[31m\x1b[39m\x1b[22m",
		mixinString)
}
