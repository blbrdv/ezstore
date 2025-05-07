package store

import "testing"

func TestNewBundlesList(t *testing.T) {
	bundle, _ := newBundleData("SomeApp.Some-Name123_1.2.3.4_X86_~_f1o2o3b4a5r6.BlockMap")
	list := newBundlesList(bundle)

	if len(list.values) != 1 {
		t.Fatal("Bundle list must contain 1 bundle")
	}

	value := list.values[0]

	if bundle.String() != value.String() {
		t.Fatalf("Bundle list must contain bundle equal to provided one")
	}
}

func TestInitBundlesList(t *testing.T) {
	list := initBundlesList()

	if len(list.values) != 0 {
		t.Fatal("Function must create empty bundle list")
	}
}

func TestBundlesListAppend(t *testing.T) {
	bundle, _ := newBundleData("SomeApp.Some-Name123_1.2.3.4_X86_~_f1o2o3b4a5r6.BlockMap")
	list := &bundlesList{values: bundlesArray{}}

	list.Append(bundle)

	if len(list.values) != 1 {
		t.Fatal("Bundle list must contain 1 bundle")
	}

	value := list.values[0]

	if bundle.String() != value.String() {
		t.Fatalf("Bundle list must contain bundle equal to appended one")
	}
}

func TestBundlesListAppendDuplicate(t *testing.T) {
	bundle1, _ := newBundleData("SomeApp.Some-Name123_1.2.3.4_X86_~_f1o2o3b4a5r6.BlockMap")
	bundle2, _ := newBundleData("SomeApp.Some-Name123_1.2.3.4_X86_~_f1o2o3b4a5r6.BlockMap")
	list := &bundlesList{values: bundlesArray{}}

	list.Append(bundle1)
	list.Append(bundle2)

	if len(list.values) != 1 {
		t.Fatal("Bundle list must contain 1 bundle")
	}

	value := list.values[0]

	if bundle1.String() != value.String() {
		t.Fatalf("Bundle list must contain bundle equal to appended first")
	}
}
