package types_test

import (
	"testing"

	"github.com/blbrdv/ezstore/internal/types"
)

func TestA(t *testing.T) {
	expected := types.Version{1, 0, 0, 0}
	raw := "v1"
	actual, err := types.New(raw)

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
	expected := types.Version{1, 2, 0, 0}
	raw := "1.2"
	actual, err := types.New(raw)

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
	expected := types.Version{1, 2, 3, 0}
	raw := "1.2.3"
	actual, err := types.New(raw)

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
	expected := types.Version{1, 2, 3, 4}
	raw := "1.2.3.4"
	actual, err := types.New(raw)

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
	a := &types.Version{1, 2, 3, 4}
	b := &types.Version{0, 2, 3, 4}

	actual := a.Compare(b)

	if actual != expected {
		t.Fatalf(`Incorrect comparsion, expected: %d, actual: %d`, expected, actual)
	}
}

func TestCompareRight(t *testing.T) {
	expected := -1
	a := &types.Version{0, 0, 0, 3}
	b := &types.Version{0, 0, 0, 4}

	actual := a.Compare(b)

	if actual != expected {
		t.Fatalf(`Incorrect comparsion, expected: %d, actual: %d`, expected, actual)
	}
}

func TestCompareEqual(t *testing.T) {
	expected := 0
	a := &types.Version{1, 2, 3, 4}
	b := &types.Version{1, 2, 3, 4}

	actual := a.Compare(b)

	if actual != expected {
		t.Fatalf(`Incorrect comparsion, expected: %d, actual: %d`, expected, actual)
	}
}
