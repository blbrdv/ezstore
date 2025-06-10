package main

import (
	"fmt"
	"strings"

	"github.com/magefile/mage/sh"
)

func getLastTag() (string, error) {
	result, err := sh.Output("git", "describe", "--tags", "--abbrev=0")
	if err != nil {
		return "", err
	}

	return result, nil
}

func getProductVersion() (string, error) {
	tag, err := getLastTag()
	if err != nil {
		return "", err
	}

	if strings.HasPrefix(tag, "v") {
		return tag[1:], nil
	} else {
		return tag, nil
	}
}

func getFileVersion(version string) (string, error) {
	tag, err := getLastTag()
	if err != nil {
		return "", err
	}

	result, err := sh.Output("git", "log", fmt.Sprintf("%s..HEAD", tag), "--oneline")
	if err != nil {
		return "", err
	}

	result = strings.Trim(result, getNewLine())
	result = strings.Trim(result, "")

	var patch int
	if result == "" {
		patch = 0
	} else {
		patch = strings.Count(result, "\n") + 1
	}

	return fmt.Sprintf("%s.%d", version, patch), nil
}
