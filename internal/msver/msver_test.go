package msver_test

import (
	"testing"

	. "github.com/blbrdv/ezstore/internal/msver"
)

var versionData = []struct {
	Name    string
	Version *Version
	Raw     string
}{
	{"TestVersionAOnly", &Version{A: 1}, "v1"},
	{"TestVersionAB", &Version{A: 1, B: 2}, "1.2"},
	{"TestVersionABC", &Version{A: 1, B: 2, C: 3}, "1.2.3"},
	{"TestVersionABCD", &Version{A: 1, B: 2, C: 3, D: 4}, "1.2.3.4"},
}

func TestVersion(t *testing.T) {
	for _, data := range versionData {
		t.Run(data.Name, func(t *testing.T) {
			expected := data.Version
			actual, err := NewVersion(data.Raw)

			if err != nil {
				t.Fatalf(`Can not parse version`)
			}

			expectedStr := expected.String()
			actualStr := actual.String()

			if actualStr != expectedStr {
				t.Fatalf(`Incorrect Version, expected: "%s", actual: "%s"`, expectedStr, actualStr)
			}
		})
	}
}

var versionCompareData = []struct {
	Name     string
	Expected int
	A        *Version
	B        *Version
}{
	{"TestCompareLeft", 1, &Version{A: 1, B: 2, C: 3, D: 4}, &Version{B: 2, C: 3, D: 4}},
	{"TestCompareRight", -1, &Version{D: 3}, &Version{D: 4}},
	{"TestCompareEqual", 0, &Version{A: 1, B: 2, C: 3, D: 4}, &Version{A: 1, B: 2, C: 3, D: 4}},
}

func TestCompare(t *testing.T) {
	for _, data := range versionCompareData {
		t.Run(data.Name, func(t *testing.T) {
			expected := data.Expected
			a := data.A
			b := data.B

			actual := a.Compare(b)

			if actual != expected {
				t.Fatalf(`Incorrect comparsion, expected: %d, actual: %d`, expected, actual)
			}
		})
	}
}
