package store

import (
	"fmt"
	"github.com/blbrdv/ezstore/internal/ms"
	"github.com/blbrdv/ezstore/internal/utils"
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestFile(t *testing.T) {
	bundle, err := newBundle("FooBar-v1.0_1.0.0.0_neutral_~_f1o2o3b4a5r6.appx", "https://example.com/foobar")
	if err != nil {
		t.Fatalf("Bundle not created: %s", err.Error())
	}
	expected := &file{bundle: bundle, dependencies: *newBundles()}
	actual := newFile(bundle)

	if !expected.Equal(actual) {
		t.Fatalf("Invalid file.%s%s", utils.NewLine, cmp.Diff(expected, actual))
	}

	expectedStr := expected.String()
	actualStr := actual.String()

	if expectedStr != actualStr {
		t.Fatalf("Invalid file string, expected \"%s\", actual \"%s\"", expectedStr, actualStr)
	}
}

var fileData = []struct {
	Name         string
	Dependencies []*bundle
	Count        int
}{
	{"TestEmptyDependencies", []*bundle{}, 0},
	{"TestOneDependencies", []*bundle{bundlesData1}, 1},
	{"TestDuplicatesDependencies", []*bundle{bundlesData1, bundlesData12}, 1},
	{"TestTwoDependencies", []*bundle{bundlesData1, bundlesData2}, 2},
}

func TestAddDependencyToFile(t *testing.T) {
	for _, data := range fileData {
		t.Run(data.Name, func(t *testing.T) {
			bundle, err := newBundle("FooBar-v1.0_1.0.0.0_neutral_~_f1o2o3b4a5r6.appx", "https://example.com/foobar")
			if err != nil {
				t.Fatalf("Bundle not created: %s", err.Error())
			}
			app := newFile(bundle)

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

func createFile(input string, url string) *file {
	bundle, err := newBundle(input, url)
	if err != nil {
		panic(fmt.Sprintf("Version not created: %s", err.Error()))
	}
	return newFile(bundle)
}

var filesData1 = createFile("Foo_1.0.0.0_neutral_~_b1a2r3.appx", "https://example.com/b1a2r31000")
var filesData2 = createFile("Bar_1.0.0.0_x64__b3a2z1.msix", "https://example.com/b3a2z11000")

var filesData = []struct {
	Name  string
	Data  []*file
	Count int
}{
	{"TestEmptyFiles", []*file{}, 0},
	{"TestOneFile", []*file{filesData1}, 1},
	{"TestTwoFiles", []*file{filesData1, filesData2}, 2},
}

func TestFiles(t *testing.T) {
	for _, data := range filesData {
		t.Run(data.Name, func(t *testing.T) {
			files := newFiles(data.Data...)

			actualCount := len(files.elements)

			if data.Count != actualCount {
				t.Fatalf("Invalid bundles count, expected %d, actual %d", data.Count, actualCount)
			}
		})
	}
}

func TestGetBundleFiles(t *testing.T) {
	bundle1, err := newBundle("Foo_1.0.0.0_neutral_~_b1a2r3.appx", "https://example.com/b1a2r31000")
	if err != nil {
		t.Fatalf("Bundle not created: %s", err.Error())
	}
	file1 := newFile(bundle1)
	bundle2, err := newBundle("Foo_1.0.1.0_neutral_~_b1a2r3.appx", "https://example.com/b1a2r31010")
	if err != nil {
		t.Fatalf("Bundle not created: %s", err.Error())
	}
	file2 := newFile(bundle2)
	bundle3, err := newBundle("Foo_1.0.0.0_arm64__b1a2r3.msix", "https://example.com/b1a2r31000")
	if err != nil {
		t.Fatalf("Bundle not created: %s", err.Error())
	}
	file3 := newFile(bundle3)

	files := newFiles(file3, file2, file1)

	version, err := ms.NewVersion("1.0.0.0")
	if err != nil {
		t.Fatalf("Version not created: %s", err.Error())
	}

	actual, err := files.Get(version, ms.Amd64)
	if err != nil {
		t.Fatalf("Function return no bundle: %s", err.Error())
	}

	if !file1.Equal(actual) {
		t.Fatalf("Function return wrong bundle: %s", cmp.Diff(file1, actual))
	}
}
