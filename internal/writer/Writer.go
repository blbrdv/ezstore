package writer

import (
	"github.com/pterm/pterm"
	"io"
)

type DebugWriter struct {
	io.Writer
}

func (w DebugWriter) Write(p []byte) (n int, err error) {
	pterm.Debug.Printf(string(p))
	return len(p) - 1, nil
}

type ErrorWriter struct {
	io.Writer
}

func (w ErrorWriter) Write(p []byte) (n int, err error) {
	pterm.Error.Printf(string(p))
	return len(p) - 1, nil
}
