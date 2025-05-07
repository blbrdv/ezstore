package utils

import (
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

	for i := 0; i < len(actual); i++ {
		if expected[i] != actual[i] {
			t.Fatalf(`Slices elements not equal, expected: "%s", actual: "%s", index %d`, expected[i], actual[i], i)
		}
	}
}
