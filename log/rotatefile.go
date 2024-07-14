package log

import (
	"context"

	"github.com/gookit/slog/rotatefile"
	"github.com/tsingshaner/go-pkg/util"
)

type FileConfig = rotatefile.Config
type FileClearConfig = rotatefile.CConfig

func NewFileWriter(fns ...func(*FileConfig)) (*rotatefile.Writer, error) {
	c := rotatefile.NewDefaultConfig()
	for _, fn := range fns {
		fn(c)
	}

	return c.Create()
}

func NewFilesClear(ctx context.Context, fns ...util.WithFn[FileClearConfig]) *rotatefile.FilesClear {
	fc := &rotatefile.FilesClear{}
	fc.WithConfig(util.BuildWithOpts(rotatefile.NewCConfig(), fns...))

	go fc.DaemonClean(func() {})
	go func() {
		<-ctx.Done()
		fc.StopDaemon()
	}()

	return fc
}
