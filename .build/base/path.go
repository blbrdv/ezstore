package base

import (
	"runtime"
	"strings"
)

// NormalizePath returns path with right slashes
func NormalizePath(path string) string {
	return strings.ReplaceAll(path, `\`, "/")
}

// PrettifyPath returns path with slashes depends on host OS
func PrettifyPath(path string) string {
	if runtime.GOOS == "windows" {
		return strings.Replace(path, "/", `\`, -1)
	} else {
		return strings.Replace(path, `\`, "/", -1)
	}
}

func PathJoin(parts ...string) string {
	return NormalizePath(strings.Join(parts, "/"))
}
