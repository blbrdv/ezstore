package ms

type FileInfo struct {
	Path    string
	Name    string
	Version *Version
}

func NewFileInfo(path, name string, version *Version) *FileInfo {
	return &FileInfo{Path: path, Name: name, Version: version}
}

type BundleFileInfo struct {
	*FileInfo
	dependencies []*FileInfo
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
