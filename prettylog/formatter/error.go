package formatter

import (
	"strings"

	"github.com/tsingshaner/go-pkg/color"
)

var ErrorField = color.UnsafeBlue(color.UnsafeBold("err"))

func FormatError(err string) string {
	return color.UnsafeRed(
		color.UnsafeItalic(
			strings.ReplaceAll(
				strings.TrimRight(err, "\n"),
				"\n", "\n      "),
		),
	)
}
