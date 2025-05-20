package utils

import "testing"

var pathsData = []struct {
	Name   string
	Input  []string
	Output string
}{
	{"TestNormalizeUnixPath", []string{"a/b", "c"}, "a\\b\\c"},
	{"TestNormalizeCombinedPath", []string{"a/b\\c"}, "a\\b\\c"},
	{"TestNormalizeWindowsPath", []string{"a\\b", "c"}, "a\\b\\c"},
	{"TestNormalizeEmptyPath", []string{""}, ""},
	{"TestNormalizeNoPath", []string{"/"}, "\\"},
}

func TestNormalizePath(t *testing.T) {
	for _, data := range pathsData {
		t.Run(data.Name, func(t *testing.T) {
			expected := data.Output
			actual := Join(data.Input...)

			if actual != expected {
				t.Fatalf(`Incorrect path normalization, expected: "%s", actual: "%s"`, expected, actual)
			}
		})
	}
}
