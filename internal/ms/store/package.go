package store

import (
	"fmt"
	"github.com/blbrdv/ezstore/internal/ms"
	"regexp"
	"strings"
)

type packageFamilyName struct {
	Name string
	ID   string
}

func (pfm *packageFamilyName) Equal(other *packageFamilyName) bool {
	return pfm.Name == other.Name && pfm.ID == other.ID
}

func (pfm *packageFamilyName) String() string {
	return fmt.Sprintf("%s_%s", pfm.Name, pfm.ID)
}

func newPackageFamilyName(input string) (*packageFamilyName, error) {
	packageRegexp := regexp.MustCompile(`^([0-9a-zA-Z.-]+)_([a-z0-9]+)$`)
	matches := packageRegexp.FindStringSubmatch(input)
	if len(matches) == 0 {
		return nil, fmt.Errorf("\"%s\" is not valid package family name", input)
	}

	return &packageFamilyName{Name: matches[1], ID: matches[2]}, nil
}

type pkg struct {
	*packageFamilyName
	Version *ms.Version
	Arch    string
}

func (p *pkg) Equal(other *pkg) bool {
	return p.packageFamilyName.Equal(other.packageFamilyName) &&
		p.Version.Equal(other.Version) &&
		p.Arch == other.Arch
}

func (p *pkg) String() string {
	return fmt.Sprintf("%s_%s_%s__%s", p.Name, p.Version.String(), p.Arch, p.ID)
}

func newPackage(input string) (*pkg, error) {
	packageRegexp := regexp.MustCompile(`^([0-9a-zA-Z.-]+)_([\d.]+)_([a-zA-Z0-9]+)_~?_([a-z0-9]+)$`)
	matches := packageRegexp.FindStringSubmatch(input)
	if len(matches) == 0 {
		return nil, fmt.Errorf("\"%s\" is not valid package", input)
	}

	pfm := &packageFamilyName{Name: matches[1], ID: matches[4]}
	version, err := ms.NewVersion(matches[2])
	if err != nil {
		return nil, fmt.Errorf("\"%s\" is not valid package", input)
	}

	return &pkg{
			packageFamilyName: pfm,
			Version:           version,
			Arch:              strings.ToLower(matches[3]),
		},
		nil
}
