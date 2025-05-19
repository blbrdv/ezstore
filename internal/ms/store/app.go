package store

import (
	"fmt"
	"maps"
)

type app struct {
	*pkg
	dependencies map[string]struct{}
}

func (a *app) Add(dependency string) {
	a.dependencies[dependency] = struct{}{}
}

func (a *app) Dependencies() []string {
	return ToSlice(maps.Keys(a.dependencies))
}

func (a *app) Equal(other *app) bool {
	return a.pkg.Equal(other.pkg) &&
		Equal(a.Dependencies(), other.Dependencies(), func(l, r string) bool {
			return l == r
		})
}

func (a *app) String() string {
	return fmt.Sprintf("%s %s", a.pkg.String(), PrettyString(a.Dependencies()))
}

func newApp(input string) (*app, error) {
	pkg, err := newPackage(input)
	if err != nil {
		return nil, err
	}

	return &app{pkg: pkg, dependencies: map[string]struct{}{}}, nil
}

type apps struct {
	elements map[string]*app
}

func (a *apps) Add(app *app) {
	if a.elements[app.pkg.String()] == nil {
		a.elements[app.pkg.String()] = app
	}
}

func (a *apps) Values() []*app {
	return ToSlice(maps.Values(a.elements))
}

func (a *apps) String() string {
	return PrettyString(a.Values())
}

func newApps(values ...*app) *apps {
	result := &apps{elements: map[string]*app{}}

	for _, value := range values {
		result.Add(value)
	}

	return result
}
