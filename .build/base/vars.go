package base

import (
	"fmt"
	"path"
)

const (
	NewLine       = "\n"
	GoMod         = "go.mod"
	BuildPath     = ".build"
	LocalPath     = "."
	RecursivePath = "..."
)

var ModFile = fmt.Sprintf("-modfile=%s", path.Join(BuildPath, GoMod))
