package store

import (
	"fmt"
	"testing"
)

func TestBundleInfo(t *testing.T) {
	expected := &bundleInfo{Name: "SomeApp.Some-Name123", ID: "f1o2o3b4a5r6"}
	actual, err := newBundleInfo("SomeApp.Some-Name123_f1o2o3b4a5r6")

	if err != nil {
		t.Fatalf(`Can not parse bundle info: %s`, err.Error())
	}

	expectedStr := expected.String()
	actualStr := actual.String()

	if actualStr != expectedStr {
		t.Fatalf(`Incorrect bundle info, expected: "%s", actual: "%s"`, expectedStr, actualStr)
	}
	if actual.Name != expected.Name {
		t.Fatalf(`Incorrect bundle info, expected name: "%s", actual: "%s"`, expected.Name, actual.Name)
	}
	if actual.ID != expected.ID {
		t.Fatalf(`Incorrect bundle info, expected ID: "%s", actual: "%s"`, expected.ID, actual.ID)
	}
}

var bundleInfoData = []struct {
	Name  string
	Value string
}{
	{"TestEmptyInput", ""},
	{"TestInvalidFormat", "foo bar 321"},
	{"TestInvalidSeparator", "SomeApp.Some-Name123/f1o2o3b4a5r6"},
	{"TestInvalidNameSymbols", "*/-+=_f1o2o3b4a5r6"},
	{"TestInvalidIDSymbols", "SomeApp.Some-Name123_*/-+="},
}

func TestInvalidBundleInfo(t *testing.T) {
	for _, data := range bundleInfoData {
		t.Run(data.Name, func(t *testing.T) {
			expected := fmt.Sprintf("\"%s\" is not valid bundle info", data.Value)
			result, err := newBundleInfo(data.Value)

			if err == nil {
				t.Fatalf(`Function must return error "%s", but return result "%s"`, expected, result.String())
			}
			if err.Error() != expected {
				t.Fatalf(`Incorrect error message, expected "%s", actual "%s"`, expected, err.Error())
			}
		})
	}
}
