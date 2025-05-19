package store

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestPackageFamilyName(t *testing.T) {
	familyName := "FooBar-v1.0_f1o2o3b4a5r6"
	pfm1, err := newPackageFamilyName(familyName)
	if err != nil {
		t.Fatalf("First package family name created with error: %s", err.Error())
	}
	pfm2, err := newPackageFamilyName(familyName)
	if err != nil {
		t.Fatalf("Second package family name created with error: %s", err.Error())
	}

	if !pfm1.Equal(pfm2) {
		t.Fatalf("Package family names not equal: %s", cmp.Diff(pfm1, pfm2))
	}

	pfm1Str := pfm1.String()
	pfm2Str := pfm2.String()

	if pfm1Str != pfm2Str {
		t.Fatalf("Package family names strings not equal, left \"%s\", right \"%s\"", pfm1Str, pfm2Str)
	}
}

func TestPackage(t *testing.T) {
	packageStr := "FooBar-v1.0_1.0.0.0_neutral_~_f1o2o3b4a5r6"
	pkg1, err := newPackage(packageStr)
	if err != nil {
		t.Fatalf("First package created with error: %s", err.Error())
	}
	pkg2, err := newPackage(packageStr)
	if err != nil {
		t.Fatalf("Second package created with error: %s", err.Error())
	}

	if !pkg1.Equal(pkg2) {
		t.Fatalf("Packages not equal: %s", cmp.Diff(pkg1, pkg2))
	}

	pkg1Str := pkg1.String()
	pkg2Str := pkg2.String()

	if pkg1Str != pkg2Str {
		t.Fatalf("Package strings not equal, left \"%s\", right \"%s\"", pkg1Str, pkg2Str)
	}
}
