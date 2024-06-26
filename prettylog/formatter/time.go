package formatter

import (
	"time"

	"github.com/tsingshaner/go-pkg/color"
)

func Time(ts time.Time) string {
	return color.UnsafeDim(ts.Format("15:04:05.999"))
}
