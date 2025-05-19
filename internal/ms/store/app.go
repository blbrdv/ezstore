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
	if !a.pkg.Equal(other.pkg) {
		return false
	}

	otherLength := len(other.Dependencies())
	if len(a.Dependencies()) != otherLength {
		return false
	}

	for i := 0; i < otherLength; i++ {
		if a.Dependencies()[i] != other.Dependencies()[i] {
			return false
		}
	}

	return true
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
