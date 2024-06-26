package formatter

import (
	"fmt"

	"github.com/tsingshaner/go-pkg/color"
)

func Pid(pid int) string {
	return color.UnsafeMagenta(fmt.Sprintf("%d", pid))
}
