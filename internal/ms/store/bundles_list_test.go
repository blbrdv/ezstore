package store

import (
	"github.com/blbrdv/ezstore/internal/ms"
	"testing"
)

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

func TestBundlesListGetSupportedX64(t *testing.T) {
	expected := "x64"
	arch := ms.Amd64
	bundle1, _ := newBundleData("SomeApp.Some-Name123_1.2.3.4_arm64_~_f1o2o3b4a5r6.BlockMap")
	bundle2, _ := newBundleData("SomeApp.Some-Name123_1.2.3.4_arm_~_f1o2o3b4a5r6.BlockMap")
	bundle3, _ := newBundleData("SomeApp.Some-Name123_1.2.3.4_x64_~_f1o2o3b4a5r6.BlockMap")
	bundle4, _ := newBundleData("SomeApp.Some-Name123_1.2.3.4_x86_~_f1o2o3b4a5r6.BlockMap")
	bundle5, _ := newBundleData("SomeApp.Some-Name123_1.2.3.4_neutral_~_f1o2o3b4a5r6.BlockMap")
	list := &bundlesList{values: bundlesArray{}}

	list.Append(bundle1)
	list.Append(bundle2)
	list.Append(bundle3)
	list.Append(bundle4)
	list.Append(bundle5)

	result, err := list.GetSupported(arch)
	if err != nil {
		t.Fatalf(`Can not get supported architecture`)
	}

	if result.Arch != expected {
		t.Fatalf("Invalid supported architecture, expected \"%s\", actual \"%s\"", expected, result.Arch)
	}
}

func TestBundlesListGetSupportedX86(t *testing.T) {
	expected := "x86"
	arch := ms.Amd64
	bundle1, _ := newBundleData("SomeApp.Some-Name123_1.2.3.4_arm64_~_f1o2o3b4a5r6.BlockMap")
	bundle2, _ := newBundleData("SomeApp.Some-Name123_1.2.3.4_arm_~_f1o2o3b4a5r6.BlockMap")
	bundle3, _ := newBundleData("SomeApp.Some-Name123_1.2.3.4_x86_~_f1o2o3b4a5r6.BlockMap")
	bundle4, _ := newBundleData("SomeApp.Some-Name123_1.2.3.4_neutral_~_f1o2o3b4a5r6.BlockMap")
	list := &bundlesList{values: bundlesArray{}}

	list.Append(bundle1)
	list.Append(bundle2)
	list.Append(bundle3)
	list.Append(bundle4)

	result, err := list.GetSupported(arch)
	if err != nil {
		t.Fatalf(`Can not get supported architecture`)
	}

	if result.Arch != expected {
		t.Fatalf("Invalid supported architecture, expected \"%s\", actual \"%s\"", expected, result.Arch)
	}
}

func TestBundlesListGetSupportedNeutral(t *testing.T) {
	expected := "neutral"
	arch := ms.Amd64
	bundle1, _ := newBundleData("SomeApp.Some-Name123_1.2.3.4_arm64_~_f1o2o3b4a5r6.BlockMap")
	bundle2, _ := newBundleData("SomeApp.Some-Name123_1.2.3.4_arm_~_f1o2o3b4a5r6.BlockMap")
	bundle3, _ := newBundleData("SomeApp.Some-Name123_1.2.3.4_neutral_~_f1o2o3b4a5r6.BlockMap")
	list := &bundlesList{values: bundlesArray{}}

	list.Append(bundle1)
	list.Append(bundle2)
	list.Append(bundle3)

	result, err := list.GetSupported(arch)
	if err != nil {
		t.Fatalf(`Can not get supported architecture`)
	}

	if result.Arch != expected {
		t.Fatalf("Invalid supported architecture, expected \"%s\", actual \"%s\"", expected, result.Arch)
	}
}
