package prettylog

import (
	"bufio"
	"fmt"
	"os"

	"github.com/bytedance/sonic"
	"github.com/tsingshaner/go-pkg/prettylog/formatter"
	"github.com/tsingshaner/go-pkg/util"
)

type ReaderOptions struct {
	Formatter    func(formatter.Data, []byte) string
	ErrorHandler func([]byte, error)
}

func JSONReader(fns ...util.WithFn[ReaderOptions]) {
	opts := util.BuildWithOpts(&ReaderOptions{}, fns...)
	if opts.Formatter == nil {
		opts.Formatter = func(data formatter.Data, log []byte) string {
			return fmt.Sprintf("%+v", data)
		}
	}
	if opts.ErrorHandler == nil {
		opts.ErrorHandler = func(log []byte, _ error) {
			fmt.Println(log)
		}
	}

	var o map[string]any
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		log := scanner.Bytes()

		if err := sonic.ConfigFastest.Unmarshal(log, &o); err != nil {
			opts.ErrorHandler(log, err)
		} else {
			fmt.Println(opts.Formatter(o, log))
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "[prettylog] unknown scanner error: (%s)\n", err)
	}
}
