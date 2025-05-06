package ms

import (
	"fmt"
	"regexp"
	"strconv"
)

// Version type represents 4-digit version used in Windows and MS Store.
type Version struct {
	Major    int64
	Minor    int64
	Build    int64
	Revision int64
}

// NewVersion returns [Version] from string representing SemVer.
func NewVersion(input string) (*Version, error) {
	semverRegexp := regexp.MustCompile(`^v?(\d+)(?:\.(\d+))?(?:\.(\d+))?(?:\.(\d+))?$`)
	matches := semverRegexp.FindStringSubmatch(input)

	if len(matches) == 0 {
		return nil, fmt.Errorf("\"%s\" is not a valid version", input)
	}

	a, err := strconv.ParseInt(matches[1], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("can not convert \"%s\" to int64: %s", matches[1], err.Error())
	}

	b, err := parse(matches[2])
	if err != nil {
		return nil, fmt.Errorf("can not convert \"%s\" to int64: %s", matches[2], err.Error())
	}

	c, err := parse(matches[3])
	if err != nil {
		return nil, fmt.Errorf("can not convert \"%s\" to int64: %s", matches[3], err.Error())
	}

	d, err := parse(matches[4])
	if err != nil {
		return nil, fmt.Errorf("can not convert \"%s\" to int64: %s", matches[4], err.Error())
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
func (v *Version) String() string {
	return fmt.Sprintf("v%d.%d.%d.%d", v.Major, v.Minor, v.Build, v.Revision)
}

// Compare two versions.
func (v *Version) Compare(other *Version) int {
	return recursiveCompare(v.Slice(), other.Slice())
}

// LessThan returns true if this [Version] less than other [Version].
func (v *Version) LessThan(other *Version) bool {
	return v.Compare(other) < 0
}

// Slice converts [Version] to array of 4 numbers.
func (v *Version) Slice() []int64 {
	return []int64{v.Major, v.Minor, v.Build, v.Revision}
}

func recursiveCompare(left []int64, right []int64) int {
	if len(left) == 0 {
		return 0
	}

	leftNumber := left[0]
	rightNumber := right[0]

	if leftNumber > rightNumber {
		return 1
	} else if leftNumber < rightNumber {
		return -1
	}

	return recursiveCompare(left[1:], right[1:])
}
