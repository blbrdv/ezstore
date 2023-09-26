package types

import (
	"fmt"
	"regexp"
	"strconv"
)

// 4-digit semver used in MS Store
type Version struct {
	A int64
	B int64
	C int64
	D int64
}

type Versions []Version

func New(version string) (*Version, error) {
	r, err := regexp.Compile(`^v?(\d+)(?:\.(\d+))?(?:\.(\d+))?(?:\.(\d+))?$`)

	if err != nil {
		return nil, err
	}

	f := r.FindStringSubmatch(version)

	a, err := strconv.ParseInt(f[1], 10, 64)

	if err != nil {
		return nil, err
	}

	var b int64
	if f[2] != "" {
		value, err := strconv.ParseInt(f[2], 10, 64)

		if err != nil {
			return nil, err
		}

		b = value
	} else {
		b = 0
	}

	var c int64
	if f[3] != "" {
		value, err := strconv.ParseInt(f[3], 10, 64)

		if err != nil {
			return nil, err
		}

		c = value
	} else {
		c = 0
	}

	var d int64
	if f[4] != "" {
		value, err := strconv.ParseInt(f[4], 10, 64)

		if err != nil {
			return nil, err
		}

		d = value
	} else {
		d = 0
	}

	return &Version{a, b, c, d}, nil
}

func (v Version) String() string {
	return fmt.Sprintf("v%d.%d.%d.%d", v.A, v.B, v.C, v.D)
}

func (v Version) Compare(versionB Version) int {
	return recursiveCompare(v.Slice(), versionB.Slice())
}

func (v Version) LessThan(versionB Version) bool {
	return v.Compare(versionB) < 0
}

func (v Version) Slice() []int64 {
	return []int64{v.A, v.B, v.C, v.D}
}

func (s Versions) Len() int {
	return len(s)
}

func (s Versions) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s Versions) Less(i, j int) bool {
	return s[i].LessThan(s[j])
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
