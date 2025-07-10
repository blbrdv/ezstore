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

var (
	winresGlob   = "rsrc_windows_*.syso"
	allowedArchs = []string{
		"arm",
		"arm64",
		"386",
		"amd64",
	}
	archsList = strings.Join(allowedArchs, ",")
	cmdPath   = base.PathJoin(base.LocalPath, "cmd")
)

func getWinresFiles(subdirs ...string) ([]string, error) {
	subdirs = append(subdirs, winresGlob)
	return filepath.Glob(base.PathJoin(subdirs...))
}

func removeWinresFiles(action *goyek.A, path string) error {
	files, err := getWinresFiles(path)
	if err != nil {
		return err
	}
	for _, file := range files {
		err := base.Remove(action, file)
		if err != nil {
			return err
		}
	}

	return nil
}

var clean = goyek.Define(goyek.Task{
	Name:  "clean",
	Usage: "Removes winres files, build and pack directories.",
	Action: func(action *goyek.A) {
		err := base.Remove(action, *buildPath, *packPath)
		if err != nil {
			action.Fatal(err)
		}

		err = removeWinresFiles(action, cmdPath)
		if err != nil {
			action.Fatal(err)
		}
	},
})

var build = goyek.Define(goyek.Task{
	Name:  "build",
	Usage: "Build project.",
	Action: func(action *goyek.A) {
		exists := true
		info, err := os.Stat(*buildPath)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				exists = false
			} else {
				action.Fatal(err)
			}
		}
		if exists && !info.IsDir() {
			action.Fatal(fmt.Errorf("build path must be dir"))
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

			fullOutputPath := base.PathJoin(*buildPath, arch)

			output, err = base.RunGoTool(
				action,
				true,
				base.LocalPath,
				nil,
				"go-winres",
				"make",
				"--arch", arch,
				"--in", "./winres.json",
				"--product-version", productVersion,
				"--file-version", fileVersion,
			)
			if len(output) > 0 {
				action.Log(output)
			}
			if err != nil {
				action.Fatal(err)
			}

			winresFiles, err := getWinresFiles()
			if err != nil {
				action.Fatal(err)
			}

			err = base.Move(action, cmdPath, winresFiles...)
			if err != nil {
				action.Fatal(err)
			}

			output, err = base.Run(
				action,
				true,
				base.LocalPath,
				map[string]string{
					"GOARCH": arch,
				},
				"go", "build",
				fmt.Sprintf("-ldflags=-X main.version=%s", productVersion),
				"-o", base.PathJoin(fullOutputPath, "bin/ezstore.exe"),
				cmdPath,
			)
			if len(output) > 0 {
				action.Log(output)
			}
			if err != nil {
				action.Fatal(err)
			}

			err = removeWinresFiles(action, cmdPath)
			if err != nil {
				action.Fatal(err)
			}

			var files []string
			err = filepath.WalkDir(cmdPath, func(path string, d fs.DirEntry, err error) error {
				if !d.IsDir() && filepath.Ext(path) != ".go" {
					files = append(files, path)
				}
				return nil
			})
			if err != nil {
				action.Fatal(err)
			}
			if len(files) < 1 {
				action.Fatal(fmt.Errorf(`could not fild dist files in "%s"`, cmdPath))
			}

			err = base.Copy(action, fullOutputPath, files...)
			if err != nil {
				action.Fatal(err)
			}
		}
	},
})

var _ = goyek.Define(goyek.Task{
	Name:  "rebuild",
	Usage: `Runs "clean" and "build" tasks.`,
	Deps: []*goyek.DefinedTask{
		clean,
		build,
	},
})
