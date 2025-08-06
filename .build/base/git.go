package base

import (
	"fmt"
	"github.com/goyek/goyek/v2"
	"strings"
)

func getLastTag(action *goyek.A, path string) (string, error) {
	result, err := Run(action, false, path, nil, "git", "describe", "--tags", "--abbrev=0")
	if err != nil {
		return "", err
	}

	return strings.TrimSuffix(result, NewLine), nil
}

// GetProductVersion returns semantic version of project
func GetProductVersion(action *goyek.A, path string) (string, error) {
	tag, err := getLastTag(action, path)
	if err != nil {
		return "", err
	}

	if strings.HasPrefix(tag, "v") {
		return tag[1:], nil
	} else {
		return tag, nil
	}
}

// GetFileVersion return semantic version of project including revision number
func GetFileVersion(action *goyek.A, version string, path string) (string, error) {
	tag, err := getLastTag(action, path)
	if err != nil {
		return "", err
	}

	result, err := Run(action, false, path, nil, "git", "log", fmt.Sprintf("%s..HEAD", tag), "--oneline")
	if err != nil {
		return "", err
	}

	result = strings.Trim(result, NewLine)
	result = strings.Trim(result, "")

	var patch int
	if result == "" {
		patch = 0
	} else {
		patch = strings.Count(result, NewLine) + 1
	}

	return fmt.Sprintf("%s.%d", version, patch), nil
}
