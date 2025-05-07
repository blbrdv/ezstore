package store

import (
	"fmt"
	"github.com/blbrdv/ezstore/internal/ms"
	"regexp"
	"strings"
)

type bundleData struct {
	*bundleInfo

	Version *ms.Version
	Arch    string
	Format  string
	URL     string
}

func newBundleData(input string) (*bundleData, error) {
	bundleRegexp := regexp.MustCompile(`^([0-9a-zA-Z.-]+)_([\d.]+)_([a-zA-Z0-9]+)_~?_([a-z0-9]+).([a-zA-Z]+)$`)
	matches := bundleRegexp.FindStringSubmatch(input)
	if len(matches) == 0 {
		return nil, fmt.Errorf("\"%s\" is not valid bundle data", input)
	}

	info := &bundleInfo{Name: matches[1], ID: matches[4]}
	version, err := ms.NewVersion(matches[2])
	if err != nil {
		return nil, fmt.Errorf("\"%s\" is not valid bundle data", input)
	}

	return &bundleData{
			bundleInfo: info,
			Version:    version,
			Arch:       strings.ToLower(matches[3]),
			Format:     strings.ToLower(matches[5]),
		},
		nil
}

func (bd *bundleData) String() string {
	return fmt.Sprintf(
		"{ Name: %s, ID: %s, Version: %s, Architecture: %s, Format: %s }",
		bd.Name,
		bd.ID,
		bd.Version.String(),
		bd.Arch,
		bd.Format,
	)
}
