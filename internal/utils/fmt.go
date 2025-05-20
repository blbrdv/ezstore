package utils

import (
	"fmt"
	"io"
)

func Fprint(w io.Writer, input ...any) {
	_, err := fmt.Fprint(w, input...)
	if err != nil {
		panic(err.Error())
	}
}

func Fprintf(w io.Writer, format string, input ...any) {
	_, err := fmt.Fprintf(w, format, input...)
	if err != nil {
		panic(err.Error())
	}
}

func Fprintln(w io.Writer, input ...any) {
	_, err := fmt.Fprintln(w, input...)
	if err != nil {
		panic(err.Error())
	}
}
