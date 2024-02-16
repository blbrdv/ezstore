package types

type BundleData struct {
	Version *Version
	Name    string
	Url     string
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
