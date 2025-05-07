package ms_test

import (
	"fmt"
	"testing"

	. "github.com/blbrdv/ezstore/internal/ms"
)

var versionData = []struct {
	Name    string
	Version *Version
	Raw     string
}{
	{"TestVersionMajor", &Version{Major: 1}, "v1"},
	{"TestVersionMajorMinor", &Version{Major: 1, Minor: 2}, "1.2"},
	{"TestVersionMajorMinorBuild", &Version{Major: 1, Minor: 2, Build: 3}, "1.2.3"},
	{"TestVersionMajorMinorBuildRevision", &Version{Major: 1, Minor: 2, Build: 3, Revision: 4}, "1.2.3.4"},
}

func TestVersion(t *testing.T) {
	for _, data := range versionData {
		t.Run(data.Name, func(t *testing.T) {
			expected := data.Version
			actual, err := NewVersion(data.Raw)

			if err != nil {
				t.Fatalf(`Can not parse version: %s`, err.Error())
			}

			expectedStr := expected.String()
			actualStr := actual.String()

			if actualStr != expectedStr {
				t.Fatalf(`Incorrect Version, expected: "%s", actual: "%s"`, expectedStr, actualStr)
			}
		})
	}
}

var invalidVersionData = []struct {
	Name  string
	Value string
}{
	{"TestEmptyInput", ""},
	{"TestTooManyNumbers", "1.2.3.4.5"},
	{"TestInvalidFormat", "foo bar 123"},
	{"TestInvalidInput", "vNotAVersion"},
}

func TestInvalidVersion(t *testing.T) {
	for _, data := range invalidVersionData {
		t.Run(data.Name, func(t *testing.T) {
			expected := fmt.Sprintf("\"%s\" is not a valid version", data.Value)
			result, err := NewVersion(data.Value)

			if err == nil {
				t.Fatalf(`Function must return error "%s", but return result "%s"`, expected, result.String())
			}
			if err.Error() != expected {
				t.Fatalf(`Incorrect error message, expected "%s", actual "%s"`, expected, err.Error())
			}
		})
	}
}

var versionCompareData = []struct {
	Name     string
	Expected int
	Left     *Version
	Right    *Version
}{
	{"TestCompareLeft", 1, &Version{Major: 1, Minor: 2, Build: 3, Revision: 4}, &Version{Minor: 2, Build: 3, Revision: 4}},
	{"TestCompareRight", -1, &Version{Revision: 3}, &Version{Revision: 4}},
	{"TestCompareEqual", 0, &Version{Major: 1, Minor: 2, Build: 3, Revision: 4}, &Version{Major: 1, Minor: 2, Build: 3, Revision: 4}},
}

func TestCompare(t *testing.T) {
	for _, data := range versionCompareData {
		t.Run(data.Name, func(t *testing.T) {
			expected := data.Expected
			actual := data.Left.Compare(data.Right)

			if actual != expected {
				t.Fatalf(`Incorrect comparsion, expected: %d, actual: %d`, expected, actual)
			}
		})
	}
}
