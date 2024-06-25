package log

import "io"

type Writer interface {
	io.Writer
}

type writerContainer struct {
	writers []io.Writer
}

func NewWriter(writers ...io.Writer) Writer {
	return &writerContainer{writers}
}

// Write todo use slog.Handler impl
func (wc *writerContainer) Write(p []byte) (int, error) {
	type writeResult struct {
		n   int
		err error
	}

	results := make(chan writeResult, len(wc.writers))
	for _, w := range wc.writers {
		go func(w io.Writer) {
			n, err := w.Write(p)
			results <- writeResult{n, err}
		}(w)
	}

	for i := 0; i < len(wc.writers); i++ {
		result := <-results
		if result.err != nil {
			return result.n, result.err
		}
	}

	return len(p), nil
}
