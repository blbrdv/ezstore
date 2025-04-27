package cmd

import (
	"github.com/pterm/pterm"
	"io"
)

// DebugWriter is a wrapper around pterm.Debug.Writer.
type DebugWriter struct {
	io.Writer
}

// Write provided bytes array to pterm.Debug.
func (w DebugWriter) Write(p []byte) (n int, err error) {
	pterm.Debug.Print(string(p))
	return len(p) - 1, nil
}

// ErrorWriter is a wrapper around pterm.Error.Writer.
type ErrorWriter struct {
	io.Writer
}

// Write provided bytes array to pterm.Error.
func (w ErrorWriter) Write(p []byte) (n int, err error) {
	pterm.Error.Print(string(p))
	return len(p) - 1, nil
}
