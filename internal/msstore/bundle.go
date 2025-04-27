package msstore

import (
	types "github.com/blbrdv/ezstore/internal"
)

type BundleData struct {
	Version *types.Version
	Name    string
	URL     string
	Arch    string
	Format  string
}

type Bundles []BundleData

func (d Bundles) Len() int {
	return len(d)
}

func (d Bundles) Less(i, j int) bool {
	return d[i].Version.Compare(d[j].Version) < 0
}

func (d Bundles) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}
