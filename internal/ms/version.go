package ms

import (
	"fmt"
	"regexp"
	"strconv"
)

// Version type represents 4-digit version used in Windows and MS Store.
type Version struct {
	Major    uint16
	Minor    uint16
	Build    uint16
	Revision uint16
	Encoded  uint64
}

var semverRegexp = regexp.MustCompile(`^v?(\d+)(?:\.(\d+))?(?:\.(\d+))?(?:\.(\d+))?$`)

// NewVersion returns [Version] from string representing SemVer.
func NewVersion(input string) (*Version, error) {
	matches := semverRegexp.FindStringSubmatch(input)

	if len(matches) == 0 {
		return nil, fmt.Errorf("\"%s\" is not a valid version", input)
	}

	major, err := parse(matches[1])
	if err != nil {
		return nil, fmt.Errorf("can not convert \"%s\" to uint16: %s", matches[1], err.Error())
	}

	minor, err := parse(matches[2])
	if err != nil {
		return nil, fmt.Errorf("can not convert \"%s\" to uint16: %s", matches[2], err.Error())
	}

	build, err := parse(matches[3])
	if err != nil {
		return nil, fmt.Errorf("can not convert \"%s\" to uint16: %s", matches[3], err.Error())
	}

	revision, err := parse(matches[4])
	if err != nil {
		return nil, fmt.Errorf("can not convert \"%s\" to uint16: %s", matches[4], err.Error())
	}

	return &Version{major, minor, build, revision, encode(major, minor, build, revision)}, nil
}

func NewVersionFromNumber(input uint64) *Version {
	major := uint16((input >> 48) & 0xFFFF)
	minor := uint16((input >> 32) & 0xFFFF)
	build := uint16((input >> 16) & 0xFFFF)
	revision := uint16(input & 0xFFFF)

	return &Version{major, minor, build, revision, input}
}

func encode(major, minor, build, revision uint16) uint64 {
	return (uint64(major) << 48) | (uint64(minor) << 32) | (uint64(build) << 16) | uint64(revision)
}

func parse(input string) (uint16, error) {
	if input != "" {
		number, err := strconv.ParseInt(input, 10, 64)
		if err != nil {
			return 0, err
		}
		if number < 0 || number > 65535 {
			return 0, fmt.Errorf("number is out of uint16 range")
		}

		return uint16(number), nil
	} else {
		return 0, nil
	}
}

// String returns SemVer representation.
func (v *Version) String() string {
	return fmt.Sprintf("v%d.%d.%d.%d", v.Major, v.Minor, v.Build, v.Revision)
}

// Number returns uint64 representation.
func (v *Version) Number() uint64 {
	return v.Encoded
}

func (v *Version) Equal(other *Version) bool {
	return v.Number() == other.Number()
}

// LessThan returns true if this [Version] less than other [Version].
func (v *Version) LessThan(other *Version) bool {
	return v.Number() < other.Number()
}

// MoreThan returns true if this [Version] more than other [Version].
func (v *Version) MoreThan(other *Version) bool {
	return v.Number() > other.Number()
}

// LessOrEqual returns true if this [Version] less or equal than other [Version].
func (v *Version) LessOrEqual(other *Version) bool {
	return v.Number() <= other.Number()
}

// MoreOrEqual returns true if this [Version] more or equal than other [Version].
func (v *Version) MoreOrEqual(other *Version) bool {
	return v.Number() >= other.Number()
}
