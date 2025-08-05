package store

import (
	"fmt"
	"github.com/blbrdv/ezstore/internal/utils"
	"io"
	"reflect"
	"slices"
	"strings"
)

func PrettyString(slice any) string {
	tp := reflect.TypeOf(slice)
	if tp.Kind() != reflect.Slice {
		panic(fmt.Sprintf("expected slice, got '%T' type", slice))
	}

	return printCollection(tp.Elem(), reflect.ValueOf(slice))
}

var numberTypes = []reflect.Kind{
	reflect.Int,
	reflect.Int8,
	reflect.Int16,
	reflect.Int32,
	reflect.Int64,
	reflect.Uint,
	reflect.Uint8,
	reflect.Uint16,
	reflect.Uint32,
	reflect.Uint64,
	reflect.Uintptr,
	reflect.Float32,
	reflect.Float64,
	reflect.Complex64,
	reflect.Complex128,
}

var collectionTypes = []reflect.Kind{
	reflect.Array,
	reflect.Map,
	reflect.Slice,
}

var stringerType = reflect.TypeOf((*fmt.Stringer)(nil)).Elem()

func printCollection(tp reflect.Type, s reflect.Value) string {
	var sb strings.Builder

	length := s.Len()
	last := length - 1
	f := getPrintFunc(tp)

	utils.Fprint(&sb, "[")

	for i := 0; i < length; i++ {
		f(&sb, s.Index(i))
		if i != last {
			utils.Fprint(&sb, ", ")
		}
	}

	utils.Fprint(&sb, "]")

	return sb.String()
}

func getPrintFunc(tp reflect.Type) func(writer io.Writer, value reflect.Value) {
	if tp.Implements(stringerType) {
		return printf("%s")
	} else if tp.Kind() == reflect.String {
		return printf("%q")
	} else if tp.Kind() == reflect.Bool {
		return printf("%t")
	} else if slices.Contains(numberTypes, tp.Kind()) {
		return printf("%d")
	} else if slices.Contains(collectionTypes, tp.Kind()) {
		return func(writer io.Writer, value reflect.Value) {
			utils.Fprint(writer, printCollection(tp.Elem(), value))
		}
	} else if tp.Kind() == reflect.Interface {
		return func(writer io.Writer, value reflect.Value) {
			f := getPrintFunc(value.Type())
			f(writer, value)
		}
	} else if tp.Kind() == reflect.Pointer {
		return func(writer io.Writer, value reflect.Value) {
			elem := value.Elem()
			f := getPrintFunc(elem.Type())
			utils.Fprint(writer, "*")
			f(writer, elem)
		}
	} else {
		return printf("%v")
	}
}

func printf(format string) func(writer io.Writer, value reflect.Value) {
	return func(writer io.Writer, value reflect.Value) {
		utils.Fprintf(writer, format, value.Interface())
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

func Contains[T interface{ ~[]E }, E comparable](left, right T) bool {
	return slices.ContainsFunc(left, func(t E) bool {
		return slices.Contains(right, t)
	})
}
