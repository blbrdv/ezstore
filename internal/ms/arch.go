package ms

import (
	"fmt"
	"runtime"
	"strings"
)

type Architecture int

const (
	amd64 Architecture = iota + 1
	i386
	arm64
	arm
)

var names = map[string]Architecture{
	"amd64": amd64,
	"386":   i386,
	"arm64": arm64,
	"arm":   arm,
}

var compatibilities = map[Architecture][]Architecture{
	amd64: {amd64, i386},
	i386:  {i386},
	arm64: {arm64, arm},
	arm:   {arm},
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
func (a Architecture) CompatibleWith() []Architecture {
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
		panic(fmt.Errorf("\"%s\" architecture is not supported: %s", goarch, err.Error()))
	}
	return arch
}
