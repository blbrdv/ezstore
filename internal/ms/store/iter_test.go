package store

import (
	"github.com/google/go-cmp/cmp"
	"maps"
	"testing"
)

func TestToSlice(t *testing.T) {
	data := map[string]string{"a": "1", "b": "2", "c": "3"}

	expected := []string{"a", "b", "c"}
	actual := ToSlice(maps.Keys(data))

	if len(expected) != len(actual) {
		t.Fatalf(`Slices length not equal, expected: %d elements, actual: %d elements`, len(expected), len(actual))
	}

	if !Equal(expected, actual, func(l, r string) bool {
		return l == r
	}) {
		t.Fatalf("Slices elements not equal\n%s", cmp.Diff(expected, actual))
	}
}
