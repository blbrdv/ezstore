package main

import (
	"errors"
	"fmt"
	"github.com/goyek/goyek/v2"
	"io/fs"
	"main/base"
	"os"
	"path"
	"slices"
	"strings"
)

var pack = goyek.Define(goyek.Task{
	Name:  "pack",
	Usage: "Pack project's files for distribution.",
	Action: func(action *goyek.A) {
		info, err := os.Stat(*buildPath)
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
			var output string

			fullOutputPath := path.Join(*packPath, arch)

			output, err = base.Run(
				action,
				true,
				base.LocalPath,
				nil,
				"7z", "a",
				"-bso0",
				"-bd",
				"-sse",
				path.Join(fullOutputPath, fmt.Sprintf("ezstore-%s-portable.7z", arch)),
				path.Join(*buildPath, "*"),
			)
			if len(output) > 0 {
				action.Log(output)
			}
			if err != nil {
				action.Fatal(err)
			}

			// TODO: add architectures
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
