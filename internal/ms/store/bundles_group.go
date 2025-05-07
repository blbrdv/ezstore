package store

import (
	"fmt"
	"github.com/blbrdv/ezstore/internal/ms"
	"github.com/blbrdv/ezstore/internal/utils"
	"maps"
)

type bundlesByVersion map[ms.Version]*bundlesList

type bundlesGroup struct {
	values bundlesByVersion
}

func newBundlesGroup(bundle *bundleData) *bundlesGroup {
	return &bundlesGroup{values: bundlesByVersion{*bundle.Version: newBundlesList(bundle)}}
}

func initBundlesGroup() *bundlesGroup {
	return &bundlesGroup{values: bundlesByVersion{}}
}

func (bg *bundlesGroup) Add(bundle *bundleData) {
	version := *bundle.Version
	b := bg.values[version]

	if b == nil {
		bg.values[version] = newBundlesList(bundle)
	} else {
		b.Append(bundle)
	}
}

func (bg *bundlesGroup) Get(version *ms.Version) (*bundleData, error) {
	var searchVersion ms.Version
	if version == nil {
		versions := utils.ToSlice(maps.Keys(bg.values))
		searchVersion = versions[0]
		for _, key := range versions[1:] {
			if key.Compare(&searchVersion) == 1 {
				searchVersion = key
			}
		}
	} else {
		searchVersion = *version
	}

	list := bg.values[searchVersion]
	if list == nil {
		return nil, fmt.Errorf("can not get bundle by version \"%s\"", searchVersion.String())
	}

	return list.GetSupported()
}

func (bg *bundlesGroup) GetLatest() (*bundleData, error) {
	return bg.Get(nil)
}
