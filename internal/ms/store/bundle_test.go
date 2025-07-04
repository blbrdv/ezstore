package store

import (
	"fmt"
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestBundle(t *testing.T) {
	bundleStr := "FooBar-v1.0_1.0.0.0_neutral_~_f1o2o3b4a5r6.appx"
	url := "https://example.com/foobar"
	bundle1, err := newBundle(bundleStr, url)
	if err != nil {
		t.Fatalf("First bundle created with error: %s", err.Error())
	}
	bundle2, err := newBundle(bundleStr, url)
	if err != nil {
		t.Fatalf("Second bundle created with error: %s", err.Error())
	}

	if !bundle1.Equal(bundle2) {
		t.Fatalf("Bundles not equal: %s", cmp.Diff(bundle1, bundle2))
	}

	bundle1Str := bundle1.String()
	bundle2Str := bundle2.String()

	if bundle1Str != bundle2Str {
		t.Fatalf("Bundle strings not equal, left \"%s\", right \"%s\"", bundle1Str, bundle2Str)
	}
}

func createBundle(input string, url string) *bundle {
	bundle, err := newBundle(input, url)
	if err != nil {
		panic(fmt.Sprintf("Bundle not created: %s", err.Error()))
	}
	return bundle
}

var bundlesData1 = createBundle("Foo_1.0.0.0_neutral_~_b1a2r3.appx", "https://example.com/b1a2r31000")
var bundlesData12 = createBundle("Foo_1.0.0.0_neutral_~_b1a2r3.appx", "https://example.com/b1a2r31000")
var bundlesData2 = createBundle("Bar_1.0.0.0_x64__b3a2z1.msix", "https://example.com/b3a2z11000")

var bundlesData = []struct {
	Name   string
	Data   []*bundle
	Count  int
	String string
}{
	{"TestEmptyBundles", []*bundle{}, 0, "[]"},
	{"TestOneBundle", []*bundle{bundlesData1}, 1, "[Foo_v1.0.0.0_neutral__b1a2r3.appx (\"https://example.com/b1a2r31000\")]"},
	{"TestDuplicateBundles", []*bundle{bundlesData1, bundlesData12}, 1, "[Foo_v1.0.0.0_neutral__b1a2r3.appx (\"https://example.com/b1a2r31000\")]"},
	{"TestTwoBundles", []*bundle{bundlesData1, bundlesData2}, 2, "[Foo_v1.0.0.0_neutral__b1a2r3.appx (\"https://example.com/b1a2r31000\"), Bar_v1.0.0.0_x64__b3a2z1.msix (\"https://example.com/b3a2z11000\")]"},
}

func TestBundles(t *testing.T) {
	for _, data := range bundlesData {
		t.Run(data.Name, func(t *testing.T) {
			bundles := newBundles(data.Data...)

			actualCount := len(bundles.Values())

			if data.Count != actualCount {
				t.Fatalf("Invalid bundles count, expected %d, actual %d", data.Count, actualCount)
			}
		})
	}
}

func TestGetAppBundles(t *testing.T) {
	app, err := newApp("Foo_1.0.0.0_neutral_~_b1a2r3", "neutral")
	if err != nil {
		t.Fatalf("App not created: %s", err.Error())
	}
	bundle1, err := newBundle("Foo_1.0.0.0_neutral_~_b1a2r3.appx", "https://example.com/b1a2r31000")
	if err != nil {
		t.Fatalf("Bundle not created: %s", err.Error())
	}
	bundle2, err := newBundle("Bar_1.0.0.0_x64__b3a2z1.msix", "https://example.com/b3a2z11000")
	if err != nil {
		t.Fatalf("Bundle not created: %s", err.Error())
	}

	bundles := newBundles(bundle2, bundle1)

	actual, err := bundles.GetAppBundle(app)

	if err != nil {
		t.Fatalf("Function return no bundle: %s", err.Error())
	}

	if !bundle1.Equal(actual) {
		t.Fatalf("Function return wrong bundle: %s", cmp.Diff(bundle1, actual))
	}
}

func TestGetNoAppBundles(t *testing.T) {
	expectedErr := "no bundle for \"Foo_v1.0.0.0_neutral__b1a2r3\""
	app, err := newApp("Foo_1.0.0.0_neutral_~_b1a2r3", "neutral")
	if err != nil {
		t.Fatalf("App not created: %s", err.Error())
	}
	bundle1, err := newBundle("Foo_1.0.1.0_neutral_~_b1a2r3.appx", "https://example.com/b1a2r31010")
	if err != nil {
		t.Fatalf("Bundle not created: %s", err.Error())
	}
	bundle2, err := newBundle("Bar_1.0.0.0_x64__b3a2z1.msix", "https://example.com/b3a2z11000")
	if err != nil {
		t.Fatalf("Bundle not created: %s", err.Error())
	}

	bundles := newBundles(bundle2, bundle1)

	actual, err := bundles.GetAppBundle(app)

	if err == nil {
		t.Fatalf("Function must return error \"%s\", but return result \"%s\"", expectedErr, actual.String())
	}

	if expectedErr != err.Error() {
		t.Fatalf(`Incorrect error message, expected "%s", actual "%s"`, expectedErr, err.Error())
	}
}

func TestGetDepBundles(t *testing.T) {
	dep := &dependency{name: "Foo"}
	bundle1, err := newBundle("Foo_1.0.0.0_neutral_~_b1a2r3.appx", "https://example.com/b1a2r31000")
	if err != nil {
		t.Fatalf("Bundle not created: %s", err.Error())
	}
	bundle2, err := newBundle("Foo_1.0.1.0_neutral_~_b1a2r3.appx", "https://example.com/b1a2r31010")
	if err != nil {
		t.Fatalf("Bundle not created: %s", err.Error())
	}
	bundle3, err := newBundle("Bar_1.0.0.0_x64__b3a2z1.msix", "https://example.com/b3a2z11000")
	if err != nil {
		t.Fatalf("Bundle not created: %s", err.Error())
	}

	bundles := newBundles(bundle3, bundle2, bundle1)

	actual, err := bundles.GetDependency(dep, "neutral")

	if err != nil {
		t.Fatalf("Function return no bundle: %s", err.Error())
	}

	if !bundle2.Equal(actual) {
		t.Fatalf("Function return wrong bundle: %s", cmp.Diff(bundle2, actual))
	}
}

func TestGetNoDepBundles(t *testing.T) {
	expectedErr := "no bundle for \"Baz\""
	dep := &dependency{name: "Baz"}
	bundle1, err := newBundle("Foo_1.0.0.0_neutral_~_b1a2r3.appx", "https://example.com/b1a2r31000")
	if err != nil {
		t.Fatalf("Bundle not created: %s", err.Error())
	}
	bundle2, err := newBundle("Foo_1.0.1.0_neutral_~_b1a2r3.appx", "https://example.com/b1a2r31010")
	if err != nil {
		t.Fatalf("Bundle not created: %s", err.Error())
	}
	bundle3, err := newBundle("Bar_1.0.0.0_x64__b3a2z1.msix", "https://example.com/b3a2z11000")
	if err != nil {
		t.Fatalf("Bundle not created: %s", err.Error())
	}

	bundles := newBundles(bundle3, bundle2, bundle1)

	actual, err := bundles.GetDependency(dep, "neutral")

	if err == nil {
		t.Fatalf("Function must return error \"%s\", but return result \"%s\"", expectedErr, actual.String())
	}

	if expectedErr != err.Error() {
		t.Fatalf(`Incorrect error message, expected "%s", actual "%s"`, expectedErr, err.Error())
	}
}
