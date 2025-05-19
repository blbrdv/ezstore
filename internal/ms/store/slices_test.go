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

func (sd *SomeData) Equal(other *SomeData) bool {
	return sd.Value == other.Value
}

var slicesPrettyStringData = []struct {
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

func TestPrettyStringSlices(t *testing.T) {
	for _, data := range slicesPrettyStringData {
		t.Run(data.Name, func(t *testing.T) {
			expected := data.Expected
			actual := PrettyString(data.Slice)

			if expected != actual {
				t.Fatalf(`Incorrect pretty string, expected: "%s", actual: "%s"`, expected, actual)
			}
		})
	}
}

var slicesEqualData = []struct {
	Name  string
	Left  []*SomeData
	Right []*SomeData
	Equal bool
}{
	{"TestEmptySlices", []*SomeData{}, []*SomeData{}, true},
	{"TestEqualSlices", []*SomeData{{Value: "a"}, {Value: "b"}}, []*SomeData{{Value: "b"}, {Value: "a"}}, true},
	{"TestUnequalLeftSlices", []*SomeData{{Value: "a"}}, []*SomeData{{Value: "b"}, {Value: "a"}}, false},
	{"TestUnequalRightSlices", []*SomeData{{Value: "a"}, {Value: "b"}}, []*SomeData{{Value: "a"}}, false},
	{"TestUnequalSlices", []*SomeData{{Value: "a"}, {Value: "b"}}, []*SomeData{{Value: "b"}, {Value: "c"}}, false},
}

func TestEqualSlices(t *testing.T) {
	for _, data := range slicesEqualData {
		t.Run(data.Name, func(t *testing.T) {
			left := data.Left
			right := data.Right
			expected := data.Equal

			actual := Equal(left, right, func(l, r *SomeData) bool {
				return l.Equal(r)
			})

			if expected != actual {
				t.Errorf("Function return incorrect value, expected %v, actual %v", expected, actual)
			}
		})
	}
}
