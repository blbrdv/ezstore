package store

import (
	"github.com/blbrdv/ezstore/internal/ms"
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestInitBundlesMap(t *testing.T) {
	bundles := initBundlesMap()

	if len(bundles.values) != 0 {
		t.Fatal("Function must create empty bundle map")
	}
}

func TestBundlesMapAdd(t *testing.T) {
	version, _ := ms.NewVersion("v1.2.3.4")
	expected, _ := newBundleData("SomeApp.Some-Name123_1.2.3.4_X86_~_f1o2o3b4a5r6.BlockMap")
	bundles := &bundlesMap{values: bundlesByID{}}

	bundles.Add(expected)

	if len(bundles.values) != 1 {
		t.Fatal("Bundle map must contain 1 bundle group")
	}

	group := bundles.values["f1o2o3b4a5r6"]

	if len(group.values) != 1 {
		t.Fatal("Bundle map must contain bundle group with 1 bundle list")
	}

	list := group.values[*version]

	if len(list.values) != 1 {
		t.Fatal("Bundle map must contain bandle group with bundle list with 1 bundle")
	}

	actual := list.values[0]

	if !expected.Equal(actual) {
		t.Fatalf("Bundle map must contain bundle equal to provided one.\n%s", cmp.Diff(expected, actual))
	}
}

func TestBundlesMapAddDuplicate(t *testing.T) {
	version, _ := ms.NewVersion("1.2.3.4")
	bundle1, _ := newBundleData("SomeApp.Some-Name123_1.2.3.4_X86_~_f1o2o3b4a5r6.BlockMap")
	bundle2, _ := newBundleData("SomeApp.Some-Name123_1.2.3.4_X86_~_f1o2o3b4a5r6.BlockMap")
	bundles := &bundlesMap{values: bundlesByID{}}

	bundles.Add(bundle1)
	bundles.Add(bundle2)

	if len(bundles.values) != 1 {
		t.Fatal("Bundle map must contain 1 bundle group")
	}

	group := bundles.values["f1o2o3b4a5r6"]

	if len(group.values) != 1 {
		t.Fatal("Bundle map must contain bundle group with 1 bundle list")
	}

	list := group.values[*version]

	if len(list.values) != 1 {
		t.Fatal("Bundle map must contain bundle list with 1 bundle")
	}

	actual := list.values[0]

	if !bundle1.Equal(actual) {
		t.Fatalf("Bundle map must contain bundle equal to added first.\n%s", cmp.Diff(bundle1, actual))
	}
}

func TestBundlesMapAddMultipleVersions(t *testing.T) {
	id := "f1o2o3b4a5r6"
	version1, _ := ms.NewVersion("1.2.3.4")
	bundle1, _ := newBundleData("SomeApp.Some-Name123_1.2.3.4_X86_~_f1o2o3b4a5r6.BlockMap")
	version2, _ := ms.NewVersion("1.2.3.5")
	bundle2, _ := newBundleData("SomeApp.Some-Name123_1.2.3.5_arm__f1o2o3b4a5r6.BlockMap")
	bundles := &bundlesMap{values: bundlesByID{}}

	bundles.Add(bundle1)

	if len(bundles.values) != 1 {
		t.Fatal("Bundle map must contain 1 bundle group")
	}

	group1 := bundles.values[id]

	if len(group1.values) != 1 {
		t.Fatal("Bundle map must contain bundle group with 1 bundle list with ID from first bundle data")
	}

	list1 := group1.values[*version1]

	if len(list1.values) != 1 {
		t.Fatal("Bundle group must contain bundle list with first added bundle data")
	}

	bundles.Add(bundle2)

	if len(bundles.values) != 1 {
		t.Fatal("Bundle map must contain 1 bundle group")
	}

	group2 := bundles.values[id]

	if len(group2.values) != 2 {
		t.Fatal("Bundle map must contain bundle group with 2 bundle lists with ID from second bundle data")
	}

	list2 := group2.values[*version2]

	if len(list2.values) != 1 {
		t.Fatal("Bundle group must contain bundle list with second added bundle data")
	}

	value1 := list1.values[0]
	value2 := list2.values[0]

	if value1.Equal(value2) {
		t.Fatal("Bundle map must contain two different bundle data with different versions")
	}
}

func TestBundlesMapAddMultipleIDs(t *testing.T) {
	id1 := "f1o2o3b4a5r6"
	version1, _ := ms.NewVersion("1.2.3.4")
	bundle1, _ := newBundleData("SomeApp.Some-Name123_1.2.3.4_X86_~_f1o2o3b4a5r6.BlockMap")
	id2 := "r6a5b4o3o2f1"
	version2, _ := ms.NewVersion("1.2.3.4")
	bundle2, _ := newBundleData("OtherApp.Other-Name456_1.2.3.4_arm__r6a5b4o3o2f1.appx")
	bundles := &bundlesMap{values: bundlesByID{}}

	bundles.Add(bundle1)

	if len(bundles.values) != 1 {
		t.Fatal("Bundle map must contain 1 bundle group")
	}

	group1 := bundles.values[id1]

	if len(group1.values) != 1 {
		t.Fatal("Bundle map must contain bundle group with 1 bundle list with ID from first bundle data")
	}

	list1 := group1.values[*version1]

	if len(list1.values) != 1 {
		t.Fatal("Bundle map must contain bundle list with first added bundle data")
	}

	bundles.Add(bundle2)

	if len(bundles.values) != 2 {
		t.Fatal("Bundle map must contain 2 bundle groups")
	}

	group2 := bundles.values[id2]

	if len(group2.values) != 1 {
		t.Fatal("Bundle map must contain bundle group with 1 bundle list with ID from second bundle data")
	}

	list2 := group2.values[*version2]

	if len(list2.values) != 1 {
		t.Fatal("Bundle group must contain bundle list with second added bundle data")
	}

	value1 := list1.values[0]
	value2 := list2.values[0]

	if value1.Equal(value2) {
		t.Fatal("Bundle map must contain two different bundle data with different versions")
	}
}
