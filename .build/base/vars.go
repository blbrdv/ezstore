package base

import (
	"fmt"
	"path"
)

var (
	NewLine       = "\n"
	GoMod         = "go.mod"
	BuildPath     = ".build"
	ModFile       = fmt.Sprintf("-modfile=%s", path.Join(BuildPath, GoMod))
	Lockfile      = fmt.Sprintf("--lockfile=%s", GoMod)
	LocalPath     = "."
	RecursivePath = "..."
)
