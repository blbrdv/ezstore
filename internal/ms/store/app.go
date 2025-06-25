package store

import (
	"fmt"
	"github.com/blbrdv/ezstore/internal/ms"
	"maps"
)

type dependency struct {
	name string
	min  *ms.Version
	max  *ms.Version
}

type app struct {
	*pkg
	dependencies map[string]*dependency
}

func (a *app) Add(name string, min, max *ms.Version) {
	a.dependencies[name] = &dependency{name, min, max}
}

func (a *app) Dependencies() []*dependency {
	var result []*dependency
	for value := range maps.Values(a.dependencies) {
		result = append(result, value)
	}
	return result
}

func (a *app) Equal(other *app) bool {
	return a.pkg.Equal(other.pkg) &&
		Equal(a.Dependencies(), other.Dependencies(), func(l, r *dependency) bool {
			return l.name == r.name
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

	return &app{pkg: pkg, dependencies: map[string]*dependency{}}, nil
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
