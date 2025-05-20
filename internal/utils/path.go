package utils

import (
	"path"
	"strings"
)

// Join joins multiple parts to one path and normalize it to Windows format.
func Join(elem ...string) string {
	result := path.Join(elem...)
	return strings.Replace(result, "/", "\\", -1)
}
