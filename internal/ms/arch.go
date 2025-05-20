package ms

import (
	"fmt"
	"runtime"
	"strings"
)

type Architecture int

const (
	Amd64 Architecture = iota + 1
	I386
	Arm64
	Arm
)

var names = map[string]Architecture{
	"x64":   Amd64, // ms
	"amd64": Amd64, // goarch
	"x86":   I386,  // ms
	"386":   I386,  // goarch
	"arm64": Arm64, // goarch, ms
	"arm":   Arm,   // goarch, ms
}

var compatibilities = map[Architecture][]string{
	Amd64: {"x64", "x86", "neutral"},
	I386:  {"x86", "neutral"},
	Arm64: {"arm64", "arm", "neutral"},
	Arm:   {"arm", "neutral"},
}

// String returns corresponding MS Store and Windows OS literal.
func (a Architecture) String() string {
	for name := range names {
		if names[name] == a {
			return name
		}
	}

	return "unknown"
}

// CompatibleWith returns a list of architectures with whom this architecture is compatible with.
func (a Architecture) CompatibleWith() []string {
	return compatibilities[a]
}

// NewArchitecture return [Architecture] from input string or error if invalid format.
func NewArchitecture(input string) (Architecture, error) {
	input = strings.TrimSpace(input)
	input = strings.ToLower(input)

	if input == "" {
		return 0, fmt.Errorf("value can not be empty")
	}

	arch := names[input]
	if arch == 0 {
		return arch, fmt.Errorf("\"%s\" is unknown architecture", input)
	}

	return arch, nil
}

var Arch = getCurrentArchitecture()

func getCurrentArchitecture() Architecture {
	goarch := runtime.GOARCH
	arch, err := NewArchitecture(goarch)
	if err != nil {
		panic(fmt.Sprintf("\"%s\" architecture is not supported: %s", goarch, err.Error()))
	}
	return arch
}
