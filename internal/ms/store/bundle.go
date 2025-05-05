package store

import (
	"fmt"
	"github.com/blbrdv/ezstore/internal/ms"
	"iter"
	"maps"
	"regexp"
	"strings"
)

func toSlice[T any](iter iter.Seq[T]) []T {
	var result []T

	for value := range iter {
		result = append(result, value)
	}

	return result
}

type bundleInfo struct {
	Name string
	ID   string
}

func newBundleInfo(input string) (*bundleInfo, error) {
	bundleRegexp := regexp.MustCompile(`^([0-9a-zA-Z.-]+)_([a-z0-9]+)$`)
	matches := bundleRegexp.FindStringSubmatch(input)
	if len(matches) == 0 {
		return nil, fmt.Errorf("%s is not valid bundle info", input)
	}

	return &bundleInfo{Name: matches[1], ID: matches[2]}, nil
}

type bundleData struct {
	*bundleInfo

	Version *ms.Version
	Arch    string
	Format  string
	URL     string
}

func newBundleData(input string) (*bundleData, error) {
	bundleRegexp := regexp.MustCompile(`^([0-9a-zA-Z.-]+)_([\d.]+)_([a-z0-9]+)_~?_([a-z0-9]+).([a-zA-Z]+)`)
	matches := bundleRegexp.FindStringSubmatch(input)
	if len(matches) == 0 {
		return nil, fmt.Errorf("%s is not valid bundle data", input)
	}

	info := &bundleInfo{Name: matches[1], ID: matches[4]}
	version, err := ms.NewVersion(matches[2])
	if err != nil {
		return nil, err
	}

	return &bundleData{
			bundleInfo: info,
			Version:    version,
			Arch:       strings.ToLower(matches[3]),
			Format:     strings.ToLower(matches[5]),
		},
		nil
}

func (bd *bundleData) String() string {
	return fmt.Sprintf(
		"{ %s, Name: %s, Version: %s, Architecture: %s, Format: %s }",
		bd.ID,
		bd.Name,
		bd.Version.String(),
		bd.Arch,
		bd.Format,
	)
}

type bundlesList []*bundleData

type bundles struct {
	bundlesList
}

func newBundles(bundle *bundleData) *bundles {
	return &bundles{bundlesList{bundle}}
}

func initBundles() *bundles {
	return &bundles{bundlesList{}}
}

func (b *bundles) Append(bundle *bundleData) {
	for _, value := range b.bundlesList {
		if value.String() == bundle.String() {
			return
		}
	}

	b.bundlesList = append(b.bundlesList, bundle)
}

func (b *bundles) GetSupported(arch ms.Architecture) (*bundleData, error) {
	for _, supported := range arch.CompatibleWith() {
		for _, data := range b.bundlesList {
			if data.Arch == supported.String() {
				return data, nil
			}
		}
	}

	for _, data := range b.bundlesList {
		if data.Arch == "neutral" {
			return data, nil
		}
	}

	return nil, fmt.Errorf("%s architecture is not supported by this app", arch.String())
}

type bundlesByVersion map[ms.Version]*bundles

type bundlesGroup struct {
	bundlesByVersion
}

func newBundlesGroup(bundle *bundleData) *bundlesGroup {
	return &bundlesGroup{bundlesByVersion{*bundle.Version: newBundles(bundle)}}
}

func initBundlesGroup() *bundlesGroup {
	return &bundlesGroup{bundlesByVersion{}}
}

func (bg *bundlesGroup) Add(bundle *bundleData) {
	version := *bundle.Version
	b := bg.bundlesByVersion[version]

	if b == nil {
		bg.bundlesByVersion[version] = newBundles(bundle)
	} else {
		b.Append(bundle)
	}
}

func (bg *bundlesGroup) Get(version *ms.Version, arch ms.Architecture) (*bundleData, error) {
	var searchVersion ms.Version
	if version == nil {
		versions := toSlice(maps.Keys(bg.bundlesByVersion))
		searchVersion = versions[0]
		for _, key := range versions[1:] {
			if key.Compare(&searchVersion) == 1 {
				searchVersion = key
			}
		}
	} else {
		searchVersion = *version
	}

	list := bg.bundlesByVersion[searchVersion]
	if list == nil {
		return nil, fmt.Errorf("can not get bundle by version %s", searchVersion.String())
	}

	return list.GetSupported(arch)
}

func (bg *bundlesGroup) GetLatest(arch ms.Architecture) (*bundleData, error) {
	return bg.Get(nil, arch)
}

type bundlesByID map[string]*bundlesGroup

type bundlesMap struct {
	bundlesByID
}

func (bm *bundlesMap) Add(bundle *bundleData) {
	group := bm.bundlesByID[bundle.ID]

	if group == nil {
		bm.bundlesByID[bundle.ID] = newBundlesGroup(bundle)
	} else {
		group.Add(bundle)
	}
}

func initBundleMap() *bundlesMap {
	return &bundlesMap{bundlesByID{}}
}
