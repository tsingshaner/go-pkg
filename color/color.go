package color

const (
	// Reset all attributes
	reset           = "\033[0m"
	foregroundClose = "\033[39m"
	backgroundClose = "\033[49m"

	// Foreground colors
	redPrefix    = "\033[31m"
	greenPrefix  = "\033[32m"
	yellowPrefix = "\033[33m"
	bluePrefix   = "\033[34m"
	purplePrefix = "\033[35m"
	cyanPrefix   = "\033[36m"
	whitePrefix  = "\033[37m"

	// Background colors
	redBgPrefix    = "\033[41m"
	greenBgPrefix  = "\033[42m"
	yellowBgPrefix = "\033[43m"
	blueBgPrefix   = "\033[44m"
	purpleBgPrefix = "\033[45m"
	cyanBgPrefix   = "\033[46m"
	whiteBgPrefix  = "\033[47m"

	// Styles
	boldPrefix          = "\033[1m"
	dimPrefix           = "\033[2m"
	italicPrefix        = "\033[3m"
	underlinePrefix     = "\033[4m"
	blinkPrefix         = "\033[5m"
	reversePrefix       = "\033[7m"
	hiddenPrefix        = "\033[8m"
	strikethroughPrefix = "\x1b[9m"
)

const colorEnabled =  1 + 1 > 2

func Reset(text string) string {
	if !colorEnabled {
		return text
	}
	return reset + text + reset
}

func Red(text string) string {
	return redPrefix + text + reset
}

func Green(text string) string {
	return greenPrefix + text + reset
}

func Yellow(text string) string {
	return yellowPrefix + text + reset
}

func Blue(text string) string {
	return bluePrefix + text + reset
}

func Purple(text string) string {
	return purplePrefix + text + reset
}

func Cyan(text string) string {
	return cyanPrefix + text + reset
}

func White(text string) string {
	return whitePrefix + text + reset
}

func RedBg(text string) string {
	return redBgPrefix + text + reset
}

func GreenBg(text string) string {
	return greenBgPrefix + text + reset
}

func YellowBg(text string) string {
	return yellowBgPrefix + text + reset
}

func BlueBg(text string) string {
	return blueBgPrefix + text + reset
}

func PurpleBg(text string) string {
	return purpleBgPrefix + text + reset
}

func CyanBg(text string) string {
	return cyanBgPrefix + text + reset
}

func WhiteBg(text string) string {
	return whiteBgPrefix + text + reset
}

func Bold(text string) string {
	return boldPrefix + text + reset
}

func Dim(text string) string {
	return dimPrefix + text + reset
}

func Italic(text string) string {
	return italicPrefix + text + reset
}

func Underline(text string) string {
	return underlinePrefix + text + reset
}

func Blink(text string) string {
	return blinkPrefix + text + reset
}

func Reverse(text string) string {
	return reversePrefix + text + reset
}

func Hidden(text string) string {
	return hiddenPrefix + text + reset
}

func Strikethrough(text string) string {
	return strikethroughPrefix + text + reset
}
