package store

import (
	"fmt"
	"github.com/blbrdv/ezstore/internal/ms"
	"maps"
	"regexp"
	"strings"
)

type bundle struct {
	*pkg
	Format string
	URL    string
}

func (b *bundle) Equal(other *bundle) bool {
	return b.pkg.Equal(other.pkg) && b.Format == other.Format && b.URL == other.URL
}

func (b *bundle) String() string {
	return fmt.Sprintf("%s_%s_%s__%s.%s (\"%s\")", b.Name, b.Version.String(), b.Arch, b.ID, b.Format, b.URL)
}

var bundleRegexp = regexp.MustCompile(`^([0-9a-zA-Z.-]+)_([\d.]+)_([a-zA-Z0-9]+)_~?_([a-z0-9]+).([a-zA-Z]+)$`)

func newBundle(input string, url string) (*bundle, error) {
	matches := bundleRegexp.FindStringSubmatch(input)
	if len(matches) == 0 {
		return nil, fmt.Errorf("\"%s\" is not valid bundle", input)
	}

	pfm := &packageFamilyName{Name: matches[1], ID: matches[4]}
	version, err := ms.NewVersion(matches[2])
	if err != nil {
		return nil, fmt.Errorf("\"%s\" is not valid bundle: %s", input, err.Error())
	}
	pkg := &pkg{
		packageFamilyName: pfm,
		Version:           version,
		Arch:              strings.ToLower(matches[3]),
	}

	return &bundle{
			pkg:    pkg,
			Format: strings.ToLower(matches[5]),
			URL:    url,
		},
		nil
}

type bundles struct {
	elements map[string]*bundle
}

func (b *bundles) Add(bundle *bundle) {
	if b.elements[bundle.String()] == nil {
		b.elements[bundle.String()] = bundle
	}
}

func (b *bundles) Values() []*bundle {
	return ToSlice(maps.Values(b.elements))
}

func (b *bundles) GetAppBundle(app *app) (*bundle, error) {
	for _, value := range b.Values() {
		if value.pkg.Equal(app.pkg) {
			return value, nil
		}
	}

	return nil, fmt.Errorf("no bundle for \"%s\"", app.pkg.String())
}

func (b *bundles) GetDependency(depData *dependency, arch string) (*bundle, error) {
	var dependencies []*bundle

	for _, value := range b.Values() {
		if value.Name == depData.name && value.Arch == arch {
			if depData.min != nil && value.Version.LessThan(depData.min) {
				continue
			}
			if depData.max != nil && value.Version.MoreThan(depData.max) {
				continue
			}
			dependencies = append(dependencies, value)
		}
	}

	length := len(dependencies)
	if length == 1 {
		return dependencies[0], nil
	} else if length > 1 {
		latest := dependencies[0]
		for i := 1; i < length; i++ {
			if dependencies[i].Version.MoreThan(latest.Version) {
				latest = dependencies[i]
			}
		}
		return latest, nil
	}

	return nil, fmt.Errorf("no bundle for \"%s\"", depData.name)
}

func (b *bundles) String() string {
	return PrettyString(b.Values())
}

func newBundles(values ...*bundle) *bundles {
	result := &bundles{elements: map[string]*bundle{}}

	for _, value := range values {
		result.Add(value)
	}

	return result
}
