package store

import (
	"testing"
)

type SomeData struct {
	Value string
}

func (sd *SomeData) String() string {
	return sd.Value
}

var slicesData = []struct {
	Name     string
	Slice    any
	Expected string
}{
	{"TestEmptyStringSlice", []string{}, "[]"},
	{"TestOneStringSlice", []string{"foo"}, "[foo]"},
	{"TestTwoStringSlice", []string{"foo", "bar"}, "[foo, bar]"},
	{"TestManyStringSlice", []string{"foo", "bar", "baz"}, "[foo, bar, baz]"},
	{"TestRepeatStringSlice", []string{"foobar", "foobar"}, "[foobar, foobar]"},

	{"TestEmptyDataSlice", []*SomeData{}, "[]"},
	{"TestOneDataSlice", []*SomeData{{Value: "foo"}}, "[foo]"},
	{"TestTwoDataSlice", []*SomeData{{Value: "foo"}, {Value: "bar"}}, "[foo, bar]"},
	{"TestManyDataSlice", []*SomeData{{Value: "foo"}, {Value: "bar"}, {Value: "baz"}}, "[foo, bar, baz]"},
	{"TestRepeatDataSlice", []*SomeData{{Value: "foobar"}, {Value: "foobar"}}, "[foobar, foobar]"},
}

func TestPrettyPrint(t *testing.T) {
	for _, data := range slicesData {
		t.Run(data.Name, func(t *testing.T) {
			expected := data.Expected
			actual := PrettyString(data.Slice)

			if expected != actual {
				t.Fatalf(`Incorrect pretty string, expected: "%s", actual: "%s"`, expected, actual)
			}
		})
	}
}
