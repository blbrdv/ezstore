package utils

import (
	"path"
	"strings"
)

func Join(elem ...string) string {
	result := path.Join(elem...)
	return strings.Replace(result, "/", "\\", -1)
}
