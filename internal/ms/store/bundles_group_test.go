package store

import (
	"fmt"
	"github.com/blbrdv/ezstore/internal/ms"
	"testing"
)

func TestNewBundlesGroup(t *testing.T) {
	version, _ := ms.NewVersion("1.2.3.4")
	bundle, _ := newBundleData("SomeApp.Some-Name123_1.2.3.4_X86_~_f1o2o3b4a5r6.BlockMap")
	group := newBundlesGroup(bundle)

	if len(group.values) != 1 {
		t.Fatal("Bundle group must contain 1 bundle list")
	}

	list := group.values[*version]

	if len(list.values) != 1 {
		t.Fatal("Bundle group must contain bundle list with 1 bundle")
	}

	value := list.values[0]

	if bundle.String() != value.String() {
		t.Fatal("Bundle group must contain bundle equal to provided one")
	}
}

func TestInitBundlesGroup(t *testing.T) {
	group := initBundlesGroup()

	if len(group.values) != 0 {
		t.Fatal("Function must create empty bundle group")
	}
}

func TestBundlesGroupAdd(t *testing.T) {
	version, _ := ms.NewVersion("1.2.3.4")
	bundle, _ := newBundleData("SomeApp.Some-Name123_1.2.3.4_X86_~_f1o2o3b4a5r6.BlockMap")
	group := &bundlesGroup{values: bundlesByVersion{}}

	group.Add(bundle)

	if len(group.values) != 1 {
		t.Fatal("Bundle group must contain 1 bundle list")
	}

	list := group.values[*version]

	if len(list.values) != 1 {
		t.Fatal("Bundle group must contain bundle list with 1 bundle")
	}

	value := list.values[0]

	if bundle.String() != value.String() {
		t.Fatal("Bundle group must contain bundle equal to appended one")
	}
}

func TestBundlesGroupAddDuplicate(t *testing.T) {
	version, _ := ms.NewVersion("1.2.3.4")
	bundle1, _ := newBundleData("SomeApp.Some-Name123_1.2.3.4_X86_~_f1o2o3b4a5r6.BlockMap")
	bundle2, _ := newBundleData("SomeApp.Some-Name123_1.2.3.4_X86_~_f1o2o3b4a5r6.BlockMap")
	group := &bundlesGroup{values: bundlesByVersion{}}

	group.Add(bundle1)
	group.Add(bundle2)

	if len(group.values) != 1 {
		t.Fatal("Bundle group must contain 1 bundle list")
	}

	list := group.values[*version]

	if len(list.values) != 1 {
		t.Fatal("Bundle group must contain bundle list with 1 bundle")
	}

	value := list.values[0]

	if bundle1.String() != value.String() {
		t.Fatal("Bundle group must contain bundle equal to added first")
	}
}

func TestBundlesGroupAddMultipleVersions(t *testing.T) {
	version1, _ := ms.NewVersion("1.2.3.4")
	bundle1, _ := newBundleData("SomeApp.Some-Name123_1.2.3.4_X86_~_f1o2o3b4a5r6.BlockMap")
	version2, _ := ms.NewVersion("1.2.3.5")
	bundle2, _ := newBundleData("SomeApp.Some-Name123_1.2.3.5_X86_~_f1o2o3b4a5r6.BlockMap")
	group := &bundlesGroup{values: bundlesByVersion{}}

	group.Add(bundle1)

	if len(group.values) != 1 {
		t.Fatal("Bundle group must contain 1 bundle list")
	}

	list1 := group.values[*version1]

	if len(list1.values) != 1 {
		t.Fatal("Bundle group must contain bundle list with first added bundle data")
	}

	group.Add(bundle2)

	if len(group.values) != 2 {
		t.Fatal("Bundle group must contain 2 bundle lists")
	}

	list2 := group.values[*version2]

	if len(list2.values) != 1 {
		t.Fatal("Bundle group must contain bundle list with 1 bundle")
	}

	value1 := list1.values[0]
	value2 := list2.values[0]

	if value1.String() == value2.String() {
		t.Fatal("Bundle group must contain two different bundle data with different versions")
	}
}

func TestBundlesGroupGet(t *testing.T) {
	version, _ := ms.NewVersion("1.2.3.4")
	bundle1, _ := newBundleData("SomeApp.Some-Name123_1.2.3.5_x64_~_f1o2o3b4a5r6.BlockMap")
	bundle2, _ := newBundleData("SomeApp.Some-Name123_1.2.3.4_x64_~_f1o2o3b4a5r6.BlockMap")
	bundle3, _ := newBundleData("SomeApp.Some-Name123_1.2.3.4_x86_~_f1o2o3b4a5r6.BlockMap")
	bundle4, _ := newBundleData("SomeApp.Some-Name123_1.2.3.3_x64_~_f1o2o3b4a5r6.BlockMap")
	group := &bundlesGroup{values: bundlesByVersion{}}

	group.Add(bundle1)
	group.Add(bundle2)
	group.Add(bundle3)
	group.Add(bundle4)

	actual, err := group.Get(version, ms.I386)
	if err != nil {
		t.Fatalf("Can not get bundle by version: %s", err.Error())
	}

	if actual.String() == bundle3.String() {
		t.Fatal("Bundle must be equal to third added bundle")
	}
}

func TestBundlesGroupGetNoVersion(t *testing.T) {
	version, _ := ms.NewVersion("1.2.3.6")
	expected := fmt.Sprintf("can not get bundle by version \"%s\"", version.String())
	bundle1, _ := newBundleData("SomeApp.Some-Name123_1.2.3.5_x64_~_f1o2o3b4a5r6.BlockMap")
	bundle2, _ := newBundleData("SomeApp.Some-Name123_1.2.3.4_x64_~_f1o2o3b4a5r6.BlockMap")
	bundle3, _ := newBundleData("SomeApp.Some-Name123_1.2.3.4_x86_~_f1o2o3b4a5r6.BlockMap")
	bundle4, _ := newBundleData("SomeApp.Some-Name123_1.2.3.3_x64_~_f1o2o3b4a5r6.BlockMap")
	group := &bundlesGroup{values: bundlesByVersion{}}

	group.Add(bundle1)
	group.Add(bundle2)
	group.Add(bundle3)
	group.Add(bundle4)

	result, err := group.Get(version, ms.I386)
	if err == nil {
		t.Fatalf(`Function must return error "%s", but return result "%s"`, expected, result.String())
	}
	if err.Error() != expected {
		t.Fatalf(`Incorrect error message, expected "%s", actual "%s"`, expected, err.Error())
	}
}

func TestBundlesGroupGetLatest(t *testing.T) {
	bundle1, _ := newBundleData("SomeApp.Some-Name123_1.2.3.5_x64_~_f1o2o3b4a5r6.BlockMap")
	bundle2, _ := newBundleData("SomeApp.Some-Name123_1.2.3.4_x64_~_f1o2o3b4a5r6.BlockMap")
	bundle3, _ := newBundleData("SomeApp.Some-Name123_1.2.3.4_x86_~_f1o2o3b4a5r6.BlockMap")
	bundle4, _ := newBundleData("SomeApp.Some-Name123_1.2.3.3_x64_~_f1o2o3b4a5r6.BlockMap")
	group := &bundlesGroup{values: bundlesByVersion{}}

	group.Add(bundle1)
	group.Add(bundle2)
	group.Add(bundle3)
	group.Add(bundle4)

	actual, err := group.GetLatest(ms.Amd64)
	if err != nil {
		t.Fatalf("Can not get bundle with latest version: %s", err.Error())
	}

	if actual.String() == bundle2.String() {
		t.Fatal("Bundle must be equal to third added bundle")
	}
}
