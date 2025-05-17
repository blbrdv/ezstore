package store

import (
	"fmt"
	"github.com/blbrdv/ezstore/internal/ms"
	"slices"
)

type file struct {
	*bundle
	dependencies bundles
}

func (f *file) Add(dependency *bundle) {
	f.dependencies.Add(dependency)
}

func (f *file) Dependencies() []*bundle {
	return f.dependencies.Values()
}

func (f *file) String() string {
	return fmt.Sprintf("%s %s", f.bundle.String(), PrettyString(f.Dependencies()))
}

func (f *file) Bundles() []*bundle {
	return slices.Concat(newBundles(f.bundle).Values(), f.Dependencies())
}

func newFile(bundle *bundle) *file {
	return &file{bundle: bundle, dependencies: *newBundles()}
}

type files struct {
	elements []*file
}

func (f *files) Add(file *file) {
	f.elements = append(f.elements, file)
}

func (f *files) Get(version *ms.Version, arch ms.Architecture) (*file, error) {
	length := len(f.elements)
	if length == 0 {
		return nil, fmt.Errorf("no file with %s %s: slice is empty", version.String(), arch)
	}

	if version == nil {
		if length == 1 {
			return f.elements[0], nil
		} else {
			latest := f.elements[0]
			for i := 1; i < length; i++ {
				if f.elements[i].Version.Compare(latest.Version) == 1 {
					latest = f.elements[i]
				}
			}
			return latest, nil
		}
	} else {
		var files []*file
		for _, file := range f.elements {
			if file.Version.Equal(version) {
				files = append(files, file)
			}
		}

		if len(files) == 0 {
			return nil, fmt.Errorf("no file with %s %s: no files with this version", version.String(), arch)
		}

		for _, supported := range arch.CompatibleWith() {
			for _, file := range files {
				if file.Arch == supported {
					return file, nil
				}
			}
		}

		return nil, fmt.Errorf("no file with %s %s: no files with supported architecture", version.String(), arch)
	}
}

func (f *files) String() string {
	return PrettyString(f.elements)
}

func newFiles(values ...*file) *files {
	result := &files{elements: []*file{}}

	for _, value := range values {
		result.Add(value)
	}

	return result
}
