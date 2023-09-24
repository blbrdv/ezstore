package msver_test

import (
	"testing"

	"github.com/blbrdv/ezstore/msver"
)

func TestA(t *testing.T) {
	expected := msver.Version{1, 0, 0, 0}
	raw := "v1"
	actual, err := msver.New(raw)

	if err != nil {
		t.Fatalf(`Can not parse version`)
	}

	expectedStr := expected.String()
	actualStr := actual.String()

	if actualStr != expectedStr {
		t.Fatalf(`Incorrect Version, expected: "%s", actual: "%s"`, expectedStr, actualStr)
	}
}

func TestB(t *testing.T) {
	expected := msver.Version{1, 2, 0, 0}
	raw := "1.2"
	actual, err := msver.New(raw)

	if err != nil {
		t.Fatalf(`Can not parse version`)
	}

	expectedStr := expected.String()
	actualStr := actual.String()

	if actualStr != expectedStr {
		t.Fatalf(`Incorrect Version, expected: "%s", actual: "%s"`, expectedStr, actualStr)
	}
}

func TestC(t *testing.T) {
	expected := msver.Version{1, 2, 3, 0}
	raw := "1.2.3"
	actual, err := msver.New(raw)

	if err != nil {
		t.Fatalf(`Can not parse version`)
	}

	expectedStr := expected.String()
	actualStr := actual.String()

	if actualStr != expectedStr {
		t.Fatalf(`Incorrect Version, expected: "%s", actual: "%s"`, expectedStr, actualStr)
	}
}

func TestD(t *testing.T) {
	expected := msver.Version{1, 2, 3, 4}
	raw := "1.2.3.4"
	actual, err := msver.New(raw)

	if err != nil {
		t.Fatalf(`Can not parse version`)
	}

	expectedStr := expected.String()
	actualStr := actual.String()

	if actualStr != expectedStr {
		t.Fatalf(`Incorrect Version, expected: "%s", actual: "%s"`, expectedStr, actualStr)
	}
}

func TestCompare(t *testing.T) {
	expected := 1
	a := msver.Version{1, 2, 3, 4}
	b := msver.Version{0, 2, 3, 4}

	actual := a.Compare(b)

	if actual != expected {
		t.Fatalf(`Incorrect comparsion, expected: %d, actual: %d`, expected, actual)
	}
}

func TestCompare2(t *testing.T) {
	expected := 1
	a := msver.Version{0, 0, 0, 4}
	b := msver.Version{0, 0, 0, 3}

	actual := a.Compare(b)

	if actual != expected {
		t.Fatalf(`Incorrect comparsion, expected: %d, actual: %d`, expected, actual)
	}
}
