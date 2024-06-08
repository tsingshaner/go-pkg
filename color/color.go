package color

import (
	"os"
	"strings"

	"github.com/mattn/go-isatty"
)

const (
	// Reset all attributes
	reset           = "\033[0m"
	foregroundClose = "\033[39m"
	backgroundClose = "\033[49m"

	// Foreground colors
	blackPrefix   = "\033[30m"
	redPrefix     = "\033[31m"
	greenPrefix   = "\033[32m"
	yellowPrefix  = "\033[33m"
	bluePrefix    = "\033[34m"
	magentaPrefix = "\033[35m"
	cyanPrefix    = "\033[36m"
	whitePrefix   = "\033[37m"
	grayPrefix    = "\033[90m"

	// Background colors
	blackBgPrefix   = "\033[40m"
	redBgPrefix     = "\033[41m"
	greenBgPrefix   = "\033[42m"
	yellowBgPrefix  = "\033[43m"
	blueBgPrefix    = "\033[44m"
	magentaBgPrefix = "\033[45m"
	cyanBgPrefix    = "\033[46m"
	whiteBgPrefix   = "\033[47m"
	grayBgPrefix    = "\033[100m"

	// Styles
	boldPrefix = "\033[1m"
	boldClose  = "\033[22m"
	dimPrefix  = "\033[2m"

	italicPrefix        = "\033[3m"
	italicClose         = "\033[23m"
	underlinePrefix     = "\033[4m"
	underlineClose      = "\033[24m"
	blinkPrefix         = "\033[5m"
	blinkClose          = "\033[25m"
	reversePrefix       = "\033[7m"
	reverseClose        = "\033[27m"
	hiddenPrefix        = "\033[8m"
	hiddenClose         = "\033[28m"
	strikethroughPrefix = "\x1b[9m"
	strikethroughClose  = "\x1b[29m"
)

var colorEnabled bool

func init() {
	ResetEnabled()
}

// IsEnabled returns true if the color output is enabled.
func IsEnabled() bool {
	return colorEnabled
}

// Enable color output.
func Enable() {
	colorEnabled = true
}

// ResetEnabled resets the color output to the default value.
func ResetEnabled() {
	colorEnabled = os.Getenv("NO_COLOR") != "" || os.Getenv("TERM") == "dumb" ||
		(!isatty.IsTerminal(os.Stdout.Fd()) && !isatty.IsCygwinTerminal(os.Stdout.Fd()))
}

// Disable color output.
func Disable() {
	colorEnabled = false
}

func format(open, text, close, replace string) string {
	if !colorEnabled {
		return text
	}

	if index := strings.Index(text, close); index != -1 {
		text = replaceClose(text, close, replace, index)
	}

	return open + text + close
}

func unsafeFormat(open, text, close string) string {
	if !colorEnabled {
		return text
	}

	return open + text + close
}

// replaceOpen replaces the opening ANSI code in the string.
func replaceClose(input, close, replace string, index int) string {
	result := ""
	cursor := 0
	for index != -1 {
		result += input[cursor:cursor+index] + replace
		cursor += index + len(close)
		index = strings.Index(input[cursor:], close)
	}
	result += input[cursor:]
	return result
}

func Reset(text string) string {
	return format(reset, text, reset, reset)
}

func UnsafeReset(text string) string {
	return unsafeFormat(reset, text, reset)
}

func Black(text string) string {
	return format(blackPrefix, text, foregroundClose, blackPrefix)
}

func UnsafeBlack(text string) string {
	return unsafeFormat(blackPrefix, text, foregroundClose)
}

func Red(text string) string {
	return format(redPrefix, text, foregroundClose, redPrefix)
}

func UnsafeRed(text string) string {
	return unsafeFormat(redPrefix, text, foregroundClose)
}

func Green(text string) string {
	return format(greenPrefix, text, foregroundClose, greenPrefix)
}

func UnsafeGreen(text string) string {
	return unsafeFormat(greenPrefix, text, foregroundClose)
}

func Yellow(text string) string {
	return format(yellowPrefix, text, foregroundClose, yellowPrefix)
}

func UnsafeYellow(text string) string {
	return unsafeFormat(yellowPrefix, text, foregroundClose)
}

func Blue(text string) string {
	return format(bluePrefix, text, foregroundClose, bluePrefix)
}

func UnsafeBlue(text string) string {
	return unsafeFormat(bluePrefix, text, foregroundClose)
}

func Magenta(text string) string {
	return format(magentaPrefix, text, foregroundClose, magentaPrefix)
}

func UnsafeMagenta(text string) string {
	return unsafeFormat(magentaPrefix, text, foregroundClose)
}

func Cyan(text string) string {
	return format(cyanPrefix, text, foregroundClose, cyanPrefix)
}

func UnsafeCyan(text string) string {
	return unsafeFormat(cyanPrefix, text, foregroundClose)
}

func White(text string) string {
	return format(whitePrefix, text, foregroundClose, whitePrefix)
}

func UnsafeWhite(text string) string {
	return unsafeFormat(whitePrefix, text, foregroundClose)
}

func Gray(text string) string {
	return format(grayPrefix, text, foregroundClose, grayPrefix)
}

func UnsafeGray(text string) string {
	return unsafeFormat(grayPrefix, text, foregroundClose)
}

func BlackBg(text string) string {
	return format(blackBgPrefix, text, backgroundClose, blackBgPrefix)
}

func UnsafeBlackBg(text string) string {
	return unsafeFormat(blackBgPrefix, text, backgroundClose)
}

func RedBg(text string) string {
	return format(redBgPrefix, text, backgroundClose, redBgPrefix)
}

func UnsafeRedBg(text string) string {
	return unsafeFormat(redBgPrefix, text, backgroundClose)
}

func GreenBg(text string) string {
	return format(greenBgPrefix, text, backgroundClose, greenBgPrefix)
}

func UnsafeGreenBg(text string) string {
	return unsafeFormat(greenBgPrefix, text, backgroundClose)
}

func YellowBg(text string) string {
	return format(yellowBgPrefix, text, backgroundClose, yellowBgPrefix)
}

func UnsafeYellowBg(text string) string {
	return unsafeFormat(yellowBgPrefix, text, backgroundClose)
}

func BlueBg(text string) string {
	return format(blueBgPrefix, text, backgroundClose, blueBgPrefix)
}

func UnsafeBlueBg(text string) string {
	return unsafeFormat(blueBgPrefix, text, backgroundClose)
}

func MagentaBg(text string) string {
	return format(magentaBgPrefix, text, backgroundClose, magentaBgPrefix)
}

func UnsafeMagentaBg(text string) string {
	return unsafeFormat(magentaBgPrefix, text, backgroundClose)
}

func CyanBg(text string) string {
	return format(cyanBgPrefix, text, backgroundClose, cyanBgPrefix)
}

func UnsafeCyanBg(text string) string {
	return unsafeFormat(cyanBgPrefix, text, backgroundClose)
}

func WhiteBg(text string) string {
	return format(whiteBgPrefix, text, backgroundClose, whiteBgPrefix)
}

func UnsafeWhiteBg(text string) string {
	return unsafeFormat(whiteBgPrefix, text, backgroundClose)
}

func GrayBg(text string) string {
	return format(grayBgPrefix, text, backgroundClose, grayBgPrefix)
}

func UnsafeGrayBg(text string) string {
	return unsafeFormat(grayBgPrefix, text, backgroundClose)
}

func Bold(text string) string {
	return format(boldPrefix, text, boldClose, "\x1b[22m\x1b[1m")
}

func UnsafeBold(text string) string {
	return unsafeFormat(boldPrefix, text, boldClose)
}

func Dim(text string) string {
	return format(dimPrefix, text, boldClose, "\x1b[22m\x1b[2m")
}

func UnsafeDim(text string) string {
	return unsafeFormat(dimPrefix, text, boldClose)
}

func Italic(text string) string {
	return format(italicPrefix, text, italicClose, italicPrefix)
}

func UnsafeItalic(text string) string {
	return unsafeFormat(italicPrefix, text, italicClose)
}

func Underline(text string) string {
	return format(underlinePrefix, text, underlineClose, underlinePrefix)
}

func UnsafeUnderline(text string) string {
	return unsafeFormat(underlinePrefix, text, underlineClose)
}

func Blink(text string) string {
	return format(blinkPrefix, text, blinkClose, blinkPrefix)
}

func UnsafeBlink(text string) string {
	return unsafeFormat(blinkPrefix, text, blinkClose)
}

func Reverse(text string) string {
	return format(reversePrefix, text, reverseClose, reversePrefix)
}

func UnsafeReverse(text string) string {
	return unsafeFormat(reversePrefix, text, reverseClose)
}

func Hidden(text string) string {
	return format(hiddenPrefix, text, hiddenClose, hiddenPrefix)
}

func UnsafeHidden(text string) string {
	return unsafeFormat(hiddenPrefix, text, hiddenClose)
}

func Strikethrough(text string) string {
	return format(strikethroughPrefix, text, strikethroughClose, strikethroughPrefix)
}

func UnsafeStrikethrough(text string) string {
	return unsafeFormat(strikethroughPrefix, text, strikethroughClose)
}
