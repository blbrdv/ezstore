package store

import (
	"fmt"
	"github.com/blbrdv/ezstore/internal/ms"
)

type bundlesArray []*bundleData

type bundlesList struct {
	values bundlesArray
}

func newBundlesList(bundle *bundleData) *bundlesList {
	return &bundlesList{values: bundlesArray{bundle}}
}

func initBundlesList() *bundlesList {
	return &bundlesList{values: bundlesArray{}}
}

func (b *bundlesList) Append(bundle *bundleData) {
	for _, value := range b.values {
		if value.String() == bundle.String() {
			return
		}
	}

	b.values = append(b.values, bundle)
}

func (b *bundlesList) GetSupported(arch ms.Architecture) (*bundleData, error) {
	for _, supported := range arch.CompatibleWith() {
		for _, data := range b.values {
			if data.Arch == supported {
				return data, nil
			}
		}
	}

	return nil, fmt.Errorf("\"%s\" architecture is not supported by this app", arch.String())
}
