package cmd

import (
	"github.com/pterm/pterm"
	"io"
)

type DebugWriter struct {
	io.Writer
}

func (w DebugWriter) Write(p []byte) (n int, err error) {
	pterm.Debug.Print(string(p))
	return len(p) - 1, nil
}

type ErrorWriter struct {
	io.Writer
}

func (w ErrorWriter) Write(p []byte) (n int, err error) {
	pterm.Error.Print(string(p))
	return len(p) - 1, nil
}
