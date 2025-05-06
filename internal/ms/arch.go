package ms

import (
	"fmt"
	"strings"
)

type Architecture int

const (
	// Amd64 represents x86-64 (also known as x64, x86_64, AMD64, and Intel 64) architecture.
	Amd64 Architecture = iota
	// Amd86 represents IA-32 (also known as "Intel Architecture, 32-bit", i386) architecture.
	Amd86
	// Arm64 represents AArch64 (also known as ARM64) architecture.
	Arm64
	// Arm86 represents ARM architecture.
	Arm86
)

var names = map[Architecture]string{
	Amd64: "x64",
	Amd86: "x86",
	Arm64: "arm64",
	Arm86: "arm",
}

var compatibilities = map[Architecture][]Architecture{
	Amd64: {Amd64, Amd86},
	Amd86: {Amd86},
	Arm64: {Arm64, Arm86},
	Arm86: {Arm86},
}

// String returns corresponding MS Store and Windows OS literal.
func (a Architecture) String() string {
	return names[a]
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
		return -1, fmt.Errorf("value can not be empty")
	}

	for arch := range names {
		if input == arch.String() {
			return arch, nil
		}
	}

	return -1, fmt.Errorf("\"%s\" is unknown architecture", input)
}
