package store

import (
	"fmt"
	"testing"
)

type ComplexStruct struct {
	Str string
	Num int
}

func (cs *ComplexStruct) String() string {
	return fmt.Sprintf("'%s' -- %d", cs.Str, cs.Num)
}

func (cs *ComplexStruct) Equal(other *ComplexStruct) bool {
	return cs.Str == other.Str
}

type SimpleStruct struct {
	Str string
	Num int
}

var slicesPrettyStringData = []struct {
	Name     string
	Slice    any
	Expected string
}{
	{"TestEmptyStringSlice", []string{}, "[]"},
	{"TestOneStringSlice", []string{"foo"}, "[\"foo\"]"},
	{"TestTwoStringSlice", []string{"foo", "bar"}, "[\"foo\", \"bar\"]"},
	{"TestManyStringSlice", []string{"foo", "bar", "baz"}, "[\"foo\", \"bar\", \"baz\"]"},
	{"TestRepeatStringSlice", []string{"foobar", "foobar"}, "[\"foobar\", \"foobar\"]"},

	{"TestEmptySimpleStructSlice", []SimpleStruct{}, "[]"},
	{"TestOneSimpleStructSlice", []SimpleStruct{{Str: "foo", Num: 1}}, "[{foo 1}]"},
	{"TestTwoSimpleStructSlice", []SimpleStruct{{Str: "foo", Num: 2}, {Str: "bar", Num: 3}}, "[{foo 2}, {bar 3}]"},
	{"TestManySimpleStructSlice", []SimpleStruct{{Str: "foo", Num: 4}, {Str: "bar", Num: 5}, {Str: "baz", Num: 6}},
		"[{foo 4}, {bar 5}, {baz 6}]"},
	{"TestRepeatSimpleStructSlice", []SimpleStruct{{Str: "foobar", Num: 7}, {Str: "foobar", Num: 7}},
		"[{foobar 7}, {foobar 7}]"},

	{"TestEmptySimpleStructPtrSlice", []*SimpleStruct{}, "[]"},
	{"TestOneSimpleStructPtrSlice", []*SimpleStruct{{Str: "foo", Num: 1}}, "[*{foo 1}]"},
	{"TestTwoSimpleStructPtrSlice", []*SimpleStruct{{Str: "foo", Num: 2}, {Str: "bar", Num: 3}},
		"[*{foo 2}, *{bar 3}]"},
	{"TestManySimpleStructPtrSlice", []*SimpleStruct{{Str: "foo", Num: 4}, {Str: "bar", Num: 5}, {Str: "baz", Num: 6}},
		"[*{foo 4}, *{bar 5}, *{baz 6}]"},
	{"TestRepeatSimpleStructPtrSlice", []*SimpleStruct{{Str: "foobar", Num: 7}, {Str: "foobar", Num: 7}},
		"[*{foobar 7}, *{foobar 7}]"},

	{"TestEmptyComplexStructSlice", []ComplexStruct{}, "[]"},
	{"TestOneComplexStructSlice", []ComplexStruct{{Str: "foo", Num: 1}}, "[{foo 1}]"},
	{"TestTwoComplexStructSlice", []ComplexStruct{{Str: "foo", Num: 2}, {Str: "bar", Num: 3}}, "[{foo 2}, {bar 3}]"},
	{"TestManyComplexStructSlice", []ComplexStruct{{Str: "foo", Num: 4}, {Str: "bar", Num: 5}, {Str: "baz", Num: 6}},
		"[{foo 4}, {bar 5}, {baz 6}]"},
	{"TestRepeatComplexStructSlice", []ComplexStruct{{Str: "foobar", Num: 7}, {Str: "foobar", Num: 7}},
		"[{foobar 7}, {foobar 7}]"},

	{"TestEmptyComplexStructPtrSlice", []*ComplexStruct{}, "[]"},
	{"TestOneComplexStructPtrSlice", []*ComplexStruct{{Str: "foo", Num: 1}}, "['foo' -- 1]"},
	{"TestTwoComplexStructPtrSlice", []*ComplexStruct{{Str: "foo", Num: 2}, {Str: "bar", Num: 3}},
		"['foo' -- 2, 'bar' -- 3]"},
	{"TestManyComplexStructPtrSlice", []*ComplexStruct{{Str: "foo", Num: 4},
		{Str: "bar", Num: 5}, {Str: "baz", Num: 6}},
		"['foo' -- 4, 'bar' -- 5, 'baz' -- 6]"},
	{"TestRepeatComplexStructPtrSlice", []*ComplexStruct{{Str: "foobar", Num: 7}, {Str: "foobar", Num: 7}},
		"['foobar' -- 7, 'foobar' -- 7]"},
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
	Left  []*ComplexStruct
	Right []*ComplexStruct
	Equal bool
}{
	{"TestEmptySlices", []*ComplexStruct{}, []*ComplexStruct{}, true},
	{"TestEqualSlices", []*ComplexStruct{{Str: "a"}, {Str: "b"}}, []*ComplexStruct{{Str: "b"}, {Str: "a"}}, true},
	{"TestUnequalLeftSlices", []*ComplexStruct{{Str: "a"}}, []*ComplexStruct{{Str: "b"}, {Str: "a"}}, false},
	{"TestUnequalRightSlices", []*ComplexStruct{{Str: "a"}, {Str: "b"}}, []*ComplexStruct{{Str: "a"}}, false},
	{"TestUnequalSlices", []*ComplexStruct{{Str: "a"}, {Str: "b"}}, []*ComplexStruct{{Str: "b"}, {Str: "c"}}, false},
}

func TestEqualSlices(t *testing.T) {
	for _, data := range slicesEqualData {
		t.Run(data.Name, func(t *testing.T) {
			left := data.Left
			right := data.Right
			expected := data.Equal

			actual := Equal(left, right, func(l, r *ComplexStruct) bool {
				return l.Equal(r)
			})

			if expected != actual {
				t.Errorf("Function return incorrect value, expected %v, actual %v", expected, actual)
			}
		})
	}
}

var containsData = []struct {
	Name     string
	Left     []string
	Right    []string
	Contains bool
}{
	{"TestContainsEqual", []string{"a", "b"}, []string{"a", "b"}, true},
	{"TestContainsLeftEqual", []string{"a", "b"}, []string{"a"}, true},
	{"TestContainsRightEqual", []string{"a"}, []string{"a", "b"}, true},
	{"TestContainsNotEqual", []string{"c", "d"}, []string{"a", "b"}, false},
}

func TestContains(t *testing.T) {
	for _, data := range containsData {
		t.Run(data.Name, func(t *testing.T) {
			left := data.Left
			right := data.Right
			expected := data.Contains

			actual := Contains(left, right)

			if actual != expected {
				t.Fatalf("expected %t, got %t", expected, actual)
			}
		})
	}
}
