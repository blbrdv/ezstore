package msstore

import (
	types "github.com/blbrdv/ezstore/internal"
)

type BundleData struct {
	Version *types.Version
	Name    string
	URL     string
	Arch    types.Architecture
	Format  string
}

type Bundles []BundleData
