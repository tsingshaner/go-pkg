package conf

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/tsingshaner/go-pkg/color"
	"github.com/tsingshaner/go-pkg/log/console"
)

// ParseArgs a function to parse the command line arguments
//
// --config  the path to the configuration file
//
// --silence silence the output of config loading
func ParseArgs() *Options {
	args := &Options{}

	flag.StringVar(&args.FilePath, "config", "config.json", "Path to the configuration file")
	flag.BoolVar(&args.Silence, "silence", false, "Silence the output of config loading")
	flag.Parse()

	if dir, err := os.Getwd(); err == nil {
		args.FilePath = filepath.Join(dir, args.FilePath)
	}

	if !args.Silence {
		showArgs(args)
	}

	return args
}

func showArgs(opts *Options) {
	console.Info(
		"will load configuration from %s",
		color.UnsafeCyan(opts.FilePath),
	)
}

func showConfig(config any) {
	sb := &strings.Builder{}
	sb.WriteString("config loaded success")
	formatStruct(sb, "    * ", config)

	console.Trace(sb.String())
}

func formatStruct(sb *strings.Builder, prefix string, obj any) {
	v := reflect.ValueOf(obj)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		sb.WriteString(fmt.Sprintf("%v", v))
		return
	}

	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		sb.WriteString(fmt.Sprintf(
			"\n%s%s %s: ",
			prefix,
			color.UnsafeCyan(t.Field(i).Name),
			color.UnsafeYellow(t.Field(i).Type.String())),
		)

		formatStruct(sb, "  "+prefix, field.Interface())
	}
}
