package log

import (
	"context"

	"github.com/gookit/slog/rotatefile"
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

func NewFilesClear(ctx context.Context, fns ...func(*FileClearConfig)) *rotatefile.FilesClear {
	c := rotatefile.NewCConfig()
	for _, fn := range fns {
		fn(c)
	}

	fc := &rotatefile.FilesClear{}
	fc.WithConfig(c)

	go fc.DaemonClean(func() {})
	go func() {
		<-ctx.Done()
		fc.StopDaemon()
	}()

	return fc
}
