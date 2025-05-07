package store

import (
	"fmt"
	"regexp"
)

type bundleInfo struct {
	Name string
	ID   string
}

func newBundleInfo(input string) (*bundleInfo, error) {
	bundleRegexp := regexp.MustCompile(`^([0-9a-zA-Z.-]+)_([a-z0-9]+)$`)
	matches := bundleRegexp.FindStringSubmatch(input)
	if len(matches) == 0 {
		return nil, fmt.Errorf("\"%s\" is not valid bundle info", input)
	}

	return &bundleInfo{Name: matches[1], ID: matches[2]}, nil
}

func (bi *bundleInfo) String() string {
	return fmt.Sprintf("{ Name: %s, ID: %s }", bi.Name, bi.ID)
}
