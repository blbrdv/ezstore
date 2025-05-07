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

func (b *bundlesList) GetSupported() (*bundleData, error) { // TODO move
	for _, supported := range ms.Arch.CompatibleWith() {
		for _, data := range b.values {
			if data.Arch == supported.String() {
				return data, nil
			}
		}
	}

	for _, data := range b.values { // TODO combine slices
		if data.Arch == "neutral" {
			return data, nil
		}
	}

	return nil, fmt.Errorf("\"%s\" architecture is not supported by this app", ms.Arch.String())
}
