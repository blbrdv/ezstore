package store

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestApp(t *testing.T) {
	pkg, _ := newPackage("Foo_1.0.0.0_neutral_~_b1a2r3")
	expected := &app{pkg: pkg, dependencies: map[string]struct{}{}}
	actual, err := newApp("Foo_1.0.0.0_neutral_~_b1a2r3")

	if err != nil {
		t.Fatalf("Function not created app: %s", err.Error())
	}

	if !expected.Equal(actual) {
		t.Fatalf("Invalid app.\n%s", cmp.Diff(expected, actual))
	}
}

var appData = []struct {
	Name         string
	Dependencies []string
	Count        int
	String       string
}{
	{"TestEmptyDependencies", []string{}, 0, "Foo_v1.0.0.0_neutral__b1a2r3 []"},
	{"TestOneDependencies", []string{"Bar"}, 1, "Foo_v1.0.0.0_neutral__b1a2r3 [Bar]"},
	{"TestDuplicatesDependencies", []string{"Bar", "Bar"}, 1, "Foo_v1.0.0.0_neutral__b1a2r3 [Bar]"},
	{"TestTwoDependencies", []string{"Bar", "Baz"}, 2, "Foo_v1.0.0.0_neutral__b1a2r3 [Bar, Baz]"},
}

func TestAddDependencyToApp(t *testing.T) {
	for _, data := range appData {
		t.Run(data.Name, func(t *testing.T) {
			app, _ := newApp("Foo_1.0.0.0_neutral_~_b1a2r3")

			for _, dep := range data.Dependencies {
				app.Add(dep)
			}

			depCount := len(app.Dependencies())

			if data.Count != depCount {
				t.Fatalf("Invalid dependency count, expected %d, actual %d", data.Count, depCount)
			}
		})
	}
}

func createApp(input string) *app {
	app, _ := newApp(input)
	return app
}

var app1 = createApp("Foo_1.0.0.0_neutral_~_b1a2r3")
var app12 = createApp("Foo_1.0.0.0_neutral_~_b1a2r3")
var app2 = createApp("Bar_1.0.0.0_x64__b3a2z1")

var appsData = []struct {
	Name   string
	Data   []*app
	Count  int
	String string
}{
	{"TestEmptyApps", []*app{}, 0, "[]"},
	{"TestOneApp", []*app{app1}, 1, "[Foo_v1.0.0.0_neutral__b1a2r3 []]"},
	{"TestDuplicateApps", []*app{app1, app12}, 1, "[Foo_v1.0.0.0_neutral__b1a2r3 []]"},
	{"TestTwoApps", []*app{app1, app2}, 2, "[Foo_v1.0.0.0_neutral__b1a2r3 [], Bar_v1.0.0.0_x64__b3a2z1 []]"},
}

func TestApps(t *testing.T) {
	for _, data := range appsData {
		t.Run(data.Name, func(t *testing.T) {
			apps := newApps(data.Data...)

			actualCount := len(apps.Values())

			if data.Count != actualCount {
				t.Fatalf("Invalid apps count, expected %d, actual %d", data.Count, actualCount)
			}
		})
	}
}
