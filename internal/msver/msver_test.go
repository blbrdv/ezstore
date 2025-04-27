package msver_test

import (
	"testing"

	. "github.com/blbrdv/ezstore/internal/msver"
)

func TestA(t *testing.T) {
	expected := Version{A: 1}
	raw := "v1"
	actual, err := NewVersion(raw)

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
	expected := Version{A: 1, B: 2}
	raw := "1.2"
	actual, err := NewVersion(raw)

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
	expected := Version{A: 1, B: 2, C: 3}
	raw := "1.2.3"
	actual, err := NewVersion(raw)

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
	expected := Version{A: 1, B: 2, C: 3, D: 4}
	raw := "1.2.3.4"
	actual, err := NewVersion(raw)

	if err != nil {
		t.Fatalf(`Can not parse version`)
	}

	expectedStr := expected.String()
	actualStr := actual.String()

	if actualStr != expectedStr {
		t.Fatalf(`Incorrect Version, expected: "%s", actual: "%s"`, expectedStr, actualStr)
	}
}

func TestCompareLeft(t *testing.T) {
	expected := 1
	a := &Version{A: 1, B: 2, C: 3, D: 4}
	b := &Version{B: 2, C: 3, D: 4}

	actual := a.Compare(b)

	if actual != expected {
		t.Fatalf(`Incorrect comparsion, expected: %d, actual: %d`, expected, actual)
	}
}

func TestCompareRight(t *testing.T) {
	expected := -1
	a := &Version{D: 3}
	b := &Version{D: 4}

	actual := a.Compare(b)

	if actual != expected {
		t.Fatalf(`Incorrect comparsion, expected: %d, actual: %d`, expected, actual)
	}
}

func TestCompareEqual(t *testing.T) {
	expected := 0
	a := &Version{A: 1, B: 2, C: 3, D: 4}
	b := &Version{A: 1, B: 2, C: 3, D: 4}

	actual := a.Compare(b)

	if actual != expected {
		t.Fatalf(`Incorrect comparsion, expected: %d, actual: %d`, expected, actual)
	}
}
