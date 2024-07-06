package prettylog

import (
	"flag"
	"strings"

	"github.com/tsingshaner/go-pkg/util/slices"
)

type CLIConfig struct {
	Presets    []string
	TimeFormat string
}

type UserConfig struct {
}

func GetCLIArgs() *CLIConfig {
	presetsStr := flag.String("presets", "default", "log presets")
	timeFormat := flag.String("time-format", "RFC3339Nano", "Time format")

	flag.Parse()

	presets := slices.Filter(
		strings.Split(*presetsStr, ","),
		func(s string) bool { return s != "" },
	)

	if len(presets) == 0 {
		presets = []string{"default"}
	}

	return &CLIConfig{
		Presets:    presets,
		TimeFormat: *timeFormat,
	}
}
