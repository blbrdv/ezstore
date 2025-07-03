package ms

import (
	"fmt"
)

type FileInfo struct {
	Path    string
	Name    string
	Version *Version
}

func NewFileInfo(path, name string, version *Version) *FileInfo {
	return &FileInfo{Path: path, Name: name, Version: version}
}

func (f *FileInfo) String() string {
	return fmt.Sprintf("%s %s '%s'", f.Name, f.Version, f.Path)
}

type BundleFileInfo struct {
	*FileInfo
	dependencies []*FileInfo
}

func (b *BundleFileInfo) String() string {
	return fmt.Sprintf("{%s %v}", b.FileInfo.String(), b.Dependencies())
}

func (b *BundleFileInfo) Dependencies() []*FileInfo {
	return b.dependencies
}

func (b *BundleFileInfo) AddDependency(file *FileInfo) {
	b.dependencies = append(b.dependencies, file)
}

func NewBundleFileInfo(file *FileInfo) *BundleFileInfo {
	return &BundleFileInfo{
		FileInfo:     file,
		dependencies: []*FileInfo{},
	}
}
