package store

import (
	"fmt"
	"io"
	"strings"
)

func fprint(w io.Writer, a ...any) {
	_, err := fmt.Fprint(w, a...)
	if err != nil {
		panic(err)
	}
}

func toString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case fmt.Stringer:
		return v.String()
	default:
		panic(fmt.Errorf("invalid value type: %T", value))
	}
}

func PrettyString[T interface{}](slice []T) string {
	var sb strings.Builder

	fprint(&sb, "[")

	length := len(slice) - 1
	for i := 0; i <= length; i++ {
		fprint(&sb, toString(slice[i]))
		if i < length {
			fprint(&sb, ", ")
		}
	}

	fprint(&sb, "]")

	return sb.String()
}
