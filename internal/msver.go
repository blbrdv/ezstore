package internal

import (
	"fmt"
	"regexp"
	"strconv"
)

// Version represents 4-digit SemVer used in Windows and MS Store.
type Version struct {
	A int64
	B int64
	C int64
	D int64
}

// NewVersion returns [Version] from string representing SemVer.
func NewVersion(input string) (*Version, error) {
	semverRegexp := regexp.MustCompile(`^v?(\d+)(?:\.(\d+))?(?:\.(\d+))?(?:\.(\d+))?$`)
	matches := semverRegexp.FindStringSubmatch(input)

	a, err := strconv.ParseInt(matches[1], 10, 64)
	if err != nil {
		return nil, err
	}

	b, err := parse(matches[2])
	if err != nil {
		return nil, err
	}

	c, err := parse(matches[3])
	if err != nil {
		return nil, err
	}

	d, err := parse(matches[4])
	if err != nil {
		return nil, err
	}

	return &Version{a, b, c, d}, nil
}

func parse(input string) (int64, error) {
	if input != "" {
		result, err := strconv.ParseInt(input, 10, 64)

		if err != nil {
			return 0, err
		}

		return result, nil
	} else {
		return 0, nil
	}
}

// String returns SemVer representation.
func (v Version) String() string {
	return fmt.Sprintf("v%d.%d.%d.%d", v.A, v.B, v.C, v.D)
}

// Compare two versions.
func (v Version) Compare(versionB *Version) int {
	return recursiveCompare(v.Slice(), versionB.Slice())
}

// LessThan returns true if this [Version] less than provided [Version].
func (v Version) LessThan(versionB *Version) bool {
	return v.Compare(versionB) < 0
}

// Slice converts [Version] to array of 4 numbers.
func (v Version) Slice() []int64 {
	return []int64{v.A, v.B, v.C, v.D}
}

func recursiveCompare(versionA []int64, versionB []int64) int {
	if len(versionA) == 0 {
		return 0
	}

	a := versionA[0]
	b := versionB[0]

	if a > b {
		return 1
	} else if a < b {
		return -1
	}

	return recursiveCompare(versionA[1:], versionB[1:])
}
