package utils

import (
	"maps"
	"testing"
)

func TestToSlice(t *testing.T) {
	data := map[string]string{"a": "1", "b": "2", "c": "3"}

	expected := []string{"a", "b", "c"}
	actual := ToSlice(func(yield func(string) bool) {
		for k := range maps.Keys(data) {
			if !yield(k) {
				return
			}
		}
	})

	if len(actual) != len(expected) {
		t.Fatalf(`Slices length not equal, expected: %d elements, actual: %d elements`, len(expected), len(actual))
	}

	for i := 0; i < len(actual); i++ {
		if actual[i] != expected[i] {
			t.Fatalf(`Slices elements not equal, expected: "%s", actual: "%s", index %d`, expected[i], actual[i], i)
		}
	}
}
