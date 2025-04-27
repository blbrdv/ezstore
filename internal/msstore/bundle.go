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
