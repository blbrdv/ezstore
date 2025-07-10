package base

import (
	"fmt"
	"io"
	"strings"
)

func Fprint(w io.Writer, a ...any) {
	_, err := fmt.Fprint(w, a...)
	if err != nil {
		panic(err.Error())
	}
}

//goland:noinspection GoUnusedExportedFunction
func Fprintln(w io.Writer, a ...any) {
	_, err := fmt.Fprintln(w, a...)
	if err != nil {
		panic(err.Error())
	}
}

//goland:noinspection GoUnusedExportedFunction
func Fprintf(w io.Writer, format string, a ...any) {
	_, err := fmt.Fprintf(w, format, a...)
	if err != nil {
		panic(err.Error())
	}
}

func Combine(output string, a any) (err error) {
	switch e := a.(type) {
	case error:
		err = e
	case string:
		err = fmt.Errorf(e)
	default:
		panic(fmt.Errorf("invalid argument type, expected error or string, got %T", a))
	}
	if len(output) == 0 {
		return err
	} else {
		if !strings.HasSuffix(output, NewLine) {
			output += NewLine
		}
		return fmt.Errorf("%s%s", output, err.Error())
	}
}
