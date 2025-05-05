package cmd

import (
	"github.com/blbrdv/ezstore/internal/log"
	"io"
)

// DebugWriter is a wrapper around log package.
type DebugWriter struct {
	io.Writer
}

// Write writes provided bytes array to os.Stdout.
func (w DebugWriter) Write(p []byte) (n int, err error) {
	log.Debug(string(p))
	return len(p) - 1, nil
}

// ErrorWriter is a wrapper around log package.
type ErrorWriter struct {
	io.Writer
}

// Write writes provided bytes array to os.Stderr.
func (w ErrorWriter) Write(p []byte) (n int, err error) {
	log.Error(string(p))
	return len(p) - 1, nil
}
