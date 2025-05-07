package store

import (
	"fmt"
	"github.com/blbrdv/ezstore/internal/ms"
	"testing"
)

func TestBundleData(t *testing.T) {
	info := &bundleInfo{Name: "SomeApp.Some-Name123", ID: "f1o2o3b4a5r6"}
	version, _ := ms.NewVersion("v1.2.3.4")
	expected := &bundleData{
		bundleInfo: info,
		Version:    version,
		Arch:       "x86",
		Format:     "blockmap",
	}
	actual, err := newBundleData("SomeApp.Some-Name123_1.2.3.4_X86_~_f1o2o3b4a5r6.BlockMap")

	if err != nil {
		t.Fatalf(`Can not parse bundle data: %s`, err.Error())
	}

	expectedStr := expected.String()
	actualStr := actual.String()

	if actualStr != expectedStr {
		t.Fatalf(`Incorrect bundle data, expected: "%s", actual: "%s"`, expectedStr, actualStr)
	}
	if actual.Name != expected.Name {
		t.Fatalf(`Incorrect bundle data, expected name: "%s", actual: "%s"`, expected.Name, actual.Name)
	}
	if actual.ID != expected.ID {
		t.Fatalf(`Incorrect bundle data, expected ID: "%s", actual: "%s"`, expected.ID, actual.ID)
	}
	if actual.Version.String() != expected.Version.String() {
		t.Fatalf(`Incorrect bundle data, expected version: "%s", actual: "%s"`, expected.Version.String(), actual.Version.String())
	}
	if actual.Arch != expected.Arch {
		t.Fatalf(`Incorrect bundle data, expected architecture: "%s", actual: "%s"`, expected.Arch, actual.Arch)
	}
	if actual.Format != expected.Format {
		t.Fatalf(`Incorrect bundle data, expected format: "%s", actual: "%s"`, expected.Format, actual.Format)
	}
}

var bundleDataData = []struct {
	Name  string
	Value string
}{
	{"TestEmptyInput", ""},
	{"TestInvalidInputFormat", "foo bar 321"},
	{"TestInvalidSeparators", "SomeApp.Some-Name123/f1o2o3b4a5r6)1.2.3.4$x86_f1o2o3b4a5r6+BlockMap"},
	{"TestInvalidNameSymbols", "*/-+=_1.2.3.4_X86_~_f1o2o3b4a5r6.BlockMap"},
	{"TestInvalidVersionSymbols", "SomeApp.Some-Name123_*/-+=_1.2.3.4_X86_~_f1o2o3b4a5r6.BlockMap"},
	{"TestInvalidArchitectureSymbols", "SomeApp.Some-Name123_1.2.3.4_*/-+=_~_f1o2o3b4a5r6.BlockMap"},
	{"TestInvalidIDSymbols", "SomeApp.Some-Name123_1.2.3.4_X86_~_*/-+=.BlockMap"},
	{"TestInvalidFormatSymbols", "SomeApp.Some-Name123_1.2.3.4_X86_~_f1o2o3b4a5r6.*/-+="},
}

func TestInvalidBundleData(t *testing.T) {
	for _, data := range bundleDataData {
		t.Run(data.Name, func(t *testing.T) {
			expected := fmt.Sprintf("\"%s\" is not valid bundle data", data.Value)
			result, err := newBundleData(data.Value)

			if err == nil {
				t.Fatalf(`Function must return error "%s", but return result "%s"`, expected, result.String())
			}
			if err.Error() != expected {
				t.Fatalf(`Incorrect error message, expected "%s", actual "%s"`, expected, err.Error())
			}
		})
	}
}
