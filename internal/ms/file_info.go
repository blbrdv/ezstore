package ms

import "fmt"

// FileInfo represents useful information about file bundle.
type FileInfo struct {
	Path    string
	Name    string
	Version *Version
}

func (fi FileInfo) String() string {
	return fmt.Sprintf("%s %s '%s'", fi.Name, fi.Version, fi.Path)
}
