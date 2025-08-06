package ms_test

import (
	"fmt"
	"github.com/blbrdv/ezstore/internal/utils"
	"github.com/google/go-cmp/cmp"
	"testing"

	. "github.com/blbrdv/ezstore/internal/ms"
)

var v1 = &Version{Major: 1, Encoded: 281474976710656}
var v12 = &Version{Major: 1, Minor: 2, Encoded: 281483566645248}
var v123 = &Version{Major: 1, Minor: 2, Build: 3, Encoded: 281483566841856}
var v1234 = &Version{Major: 1, Minor: 2, Build: 3, Revision: 4, Encoded: 281483566841860}

var versionData = []struct {
	Name    string
	Version *Version
	Raw     string
}{
	{"TestVersionMajor", v1, "v1"},
	{"TestVersionMajorMinor", v12, "1.2"},
	{"TestVersionMajorMinorBuild", v123, "1.2.3"},
	{"TestVersionMajorMinorBuildRevision", v1234, "1.2.3.4"},
}

func TestVersion(t *testing.T) {
	for _, data := range versionData {
		t.Run(data.Name, func(t *testing.T) {
			expected := data.Version
			actual, err := NewVersion(data.Raw)

			if err != nil {
				t.Fatalf(`Can not parse version: %s`, err.Error())
			}

			if !expected.Equal(actual) {
				t.Fatalf("Incorrect Version.%s%s", utils.NewLine, cmp.Diff(expected, actual))
			}

			expectedStr := expected.String()
			actualStr := actual.String()

			if actualStr != expectedStr {
				t.Fatalf(`Incorrect Version string, expected: "%s", actual: "%s"`, expectedStr, actualStr)
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

var versionDataFromNumber = []struct {
	Name    string
	Version *Version
	Raw     uint64
}{
	{"TestVersionMajor", v1, 281474976710656},
	{"TestVersionMajorMinor", v12, 281483566645248},
	{"TestVersionMajorMinorBuild", v123, 281483566841856},
	{"TestVersionMajorMinorBuildRevision", v1234, 281483566841860},
}

func TestVersionFromNumber(t *testing.T) {
	for _, data := range versionDataFromNumber {
		t.Run(data.Name, func(t *testing.T) {
			expected := data.Version
			actual := NewVersionFromNumber(data.Raw)

			if !expected.Equal(actual) {
				t.Fatalf("Incorrect Version.%s%s", utils.NewLine, cmp.Diff(expected, actual))
			}

			expectedStr := expected.String()
			actualStr := actual.String()

			if actualStr != expectedStr {
				t.Fatalf(`Incorrect Version string, expected: "%s", actual: "%s"`, expectedStr, actualStr)
			}
		})
	}
}

var versionCompareData = []struct {
	Name       string
	Comparison string
	Left       *Version
	Right      *Version
	Expected   bool
}{
	{"TestCompareEqual", "=", v12, v12, true},
	{"TestCompareNotEqual", "=", v1, v12, false},
	{"TestCompareMore", ">", v12, v1, true},
	{"TestCompareNotMore", ">", v1, v12, false},
	{"TestCompareLess", "<", v1, v12, true},
	{"TestCompareNotLess", "<", v12, v1, false},
	{"TestCompareMoreOrEqual-More", ">=", v12, v1, true},
	{"TestCompareMoreOrEqual-Equal", ">=", v12, v12, true},
	{"TestCompareNotMoreOrEqual", ">=", v1, v12, false},
	{"TestCompareLessOrEqual-More", "<=", v1, v12, true},
	{"TestCompareLessOrEqual-Equal", "<=", v12, v12, true},
	{"TestCompareNotLessOrEqual", "<=", v12, v1, false},
}

func TestCompare(t *testing.T) {
	for _, data := range versionCompareData {
		t.Run(data.Name, func(t *testing.T) {
			var actual bool
			expected := data.Expected

			switch data.Comparison {
			case "=":
				actual = data.Left.Equal(data.Right)
			case ">":
				actual = data.Left.MoreThan(data.Right)
			case "<":
				actual = data.Left.LessThan(data.Right)
			case ">=":
				actual = data.Left.MoreOrEqual(data.Right)
			case "<=":
				actual = data.Left.LessOrEqual(data.Right)
			default:
				t.Fatalf(`Invalid comparison symbol: %s`, data.Comparison)
			}

			if actual != expected {
				t.Fatalf(`Incorrect comparsion, expected: %t, actual: %t. Left: %d, right: %d`, expected, actual, data.Left.Number(), data.Right.Number())
			}
		})
	}
}
