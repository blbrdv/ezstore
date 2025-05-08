package store

type bundlesByID map[string]*bundlesGroup

type bundlesMap struct {
	values bundlesByID
}

func (bm *bundlesMap) Add(bundle *bundleData) {
	group := bm.values[bundle.ID]

	if group == nil {
		bm.values[bundle.ID] = newBundlesGroup(bundle)
	} else {
		group.Add(bundle)
	}
}

func initBundlesMap() *bundlesMap {
	return &bundlesMap{values: bundlesByID{}}
}
