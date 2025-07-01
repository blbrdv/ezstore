package store

import (
	"fmt"
	"github.com/blbrdv/ezstore/internal/ms"
)

type file struct {
	*bundle
	dependencies bundles
}

func (f *file) GetBundle() *bundle {
	return f.bundle
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

func (f *file) Equal(other *file) bool {
	return f.bundle.Equal(other.bundle) &&
		Equal(f.Dependencies(), other.Dependencies(), func(l, r *bundle) bool {
			return l.Equal(r)
		})
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

		for _, file := range files {
			if ms.IsSupported(file.Arch, arch) {
				return file, nil
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
