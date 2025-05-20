package store

import (
	"fmt"
	"github.com/blbrdv/ezstore/internal/utils"
	"reflect"
	"strings"
)

func PrettyString(slice any) string {
	val := reflect.ValueOf(slice)
	if val.Kind() != reflect.Slice {
		panic(fmt.Sprintf("expected slice, got %T", slice))
	}

	elemType := val.Type().Elem()
	if elemType.Kind() == reflect.Ptr && elemType.Elem().Kind() == reflect.Struct {
		stringerType := reflect.TypeOf((*fmt.Stringer)(nil)).Elem()
		if elemType.Implements(stringerType) {
			var sb strings.Builder

			utils.Fprint(&sb, "[")

			length := val.Len() - 1
			for i := 0; i <= length; i++ {
				v := val.Index(i).Interface().(fmt.Stringer)
				utils.Fprint(&sb, v.String())
				if i < length {
					utils.Fprint(&sb, ", ")
				}
			}

			utils.Fprint(&sb, "]")

			return sb.String()
		}
	}

	switch s := slice.(type) {
	case []string:
		var sb strings.Builder

		utils.Fprint(&sb, "[")

		length := len(s) - 1
		for i := 0; i <= length; i++ {
			utils.Fprint(&sb, s[i])
			if i < length {
				utils.Fprint(&sb, ", ")
			}
		}

		utils.Fprint(&sb, "]")

		return sb.String()
	default:
		panic(fmt.Sprintf("invalid slice type: %T", slice))
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
