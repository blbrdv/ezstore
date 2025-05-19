package store

import (
	"fmt"
	"io"
	"reflect"
	"strings"
)

func fprint(w io.Writer, a ...any) {
	_, err := fmt.Fprint(w, a...)
	if err != nil {
		panic(err)
	}
}

func PrettyString(slice any) string {
	val := reflect.ValueOf(slice)
	if val.Kind() != reflect.Slice {
		panic(fmt.Errorf("expected slice, got %T", slice))
	}

	elemType := val.Type().Elem()
	if elemType.Kind() == reflect.Ptr && elemType.Elem().Kind() == reflect.Struct {
		stringerType := reflect.TypeOf((*fmt.Stringer)(nil)).Elem()
		if elemType.Implements(stringerType) {
			var sb strings.Builder

			fprint(&sb, "[")

			length := val.Len() - 1
			for i := 0; i <= length; i++ {
				v := val.Index(i).Interface().(fmt.Stringer)
				fprint(&sb, v.String())
				if i < length {
					fprint(&sb, ", ")
				}
			}

			fprint(&sb, "]")

			return sb.String()
		}
	}

	switch s := slice.(type) {
	case []string:
		var sb strings.Builder

		fprint(&sb, "[")

		length := len(s) - 1
		for i := 0; i <= length; i++ {
			fprint(&sb, s[i])
			if i < length {
				fprint(&sb, ", ")
			}
		}

		fprint(&sb, "]")

		return sb.String()
	default:
		panic(fmt.Errorf("invalid slice type: %T", slice))
	}
}

func Equal[T any](left, right []T, f func(l, r T) bool) bool {
	if len(left) != len(right) {
		return false
	}

	for _, leftElement := range left {
		found := false
		for _, rightElement := range right {
			if f(leftElement, rightElement) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
}
