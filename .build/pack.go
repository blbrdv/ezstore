package main

import (
	"errors"
	"fmt"
	"github.com/goyek/goyek/v2"
	"io/fs"
	"main/base"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

var OSNames = map[string]string{
	"amd64": "x64",
	"386":   "x86",
	"arm64": "arm64",
	"arm":   "arm32",
}

var pack = goyek.Define(goyek.Task{
	Name:  "pack",
	Usage: "Pack project's files for distribution.",
	Action: func(action *goyek.A) {
		var output string
		var srcPath string
		if filepath.IsAbs(*buildPath) {
			srcPath = *buildPath
		} else {
			// for some reason 7z behave different for 'dir' and './dir'
			srcPath = base.PathJoin(base.LocalPath, *buildPath)
		}

		info, err := os.Stat(srcPath)
		if err != nil {
			action.Fatal(err)
		}

		exists := true
		info, err = os.Stat(*packPath)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				exists = false
			} else {
				action.Fatal(err)
			}
		}
		if exists && !info.IsDir() {
			action.Fatal(fmt.Errorf("pack path must be dir"))
		}

		targetArchs := strings.Split(*targetArch, ",")
		for _, arch := range targetArchs {
			if !slices.Contains(allowedArchs, arch) {
				action.Fatal(fmt.Errorf(`"%s" is invalid target architecture, allowed values: %v`, arch, archsList))
			}
		}

		productVersion, err := base.GetProductVersion(action, base.LocalPath)
		if err != nil {
			action.Fatal(err)
		}

		fileVersion, err := base.GetFileVersion(action, productVersion, base.LocalPath)
		if err != nil {
			action.Fatal(err)
		}

		for _, arch := range targetArchs {
			fullOutputPath := base.PathJoin(*packPath, arch)

			output, err = base.Run(
				action,
				true,
				base.LocalPath,
				nil,
				"7z", "a",
				"-bso0",
				"-bd",
				"-sse",
				base.PathJoin(fullOutputPath, fmt.Sprintf("ezstore-%s-portable.7z", OSNames[arch])),
				base.PathJoin(srcPath, arch, "*"),
			)
			if len(output) > 0 {
				action.Log(output)
			}
			if err != nil {
				action.Fatal(err)
			}
		}

		output, err = base.Run(
			action,
			true,
			base.LocalPath,
			nil,
			"iscc",
			"/Q",
			"setup.iss",
			fmt.Sprintf("/DPV=%s", productVersion),
			fmt.Sprintf("/DFV=%s", fileVersion),
		)
		if len(output) > 0 {
			action.Log(output)
		}
		if err != nil {
			action.Fatal(err)
		}
	},
})

var _ = goyek.Define(goyek.Task{
	Name:  "repack",
	Usage: `Runs "clean", "build" and "pack" tasks.`,
	Deps: []*goyek.DefinedTask{
		clean,
		build,
		pack,
	},
})
