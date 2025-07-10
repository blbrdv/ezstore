package main

import (
	"fmt"
	"github.com/goyek/goyek/v2"
	"main/base"
	"path"
	"path/filepath"
)

func getPath(value *string) string {
	if len(*value) == 0 {
		return base.LocalPath
	} else {
		return *value
	}
}

var format = goyek.Define(goyek.Task{
	Name:  "fmt",
	Usage: "Format source code.",
	Action: func(action *goyek.A) {
		targetPath := getPath(targetPath) + base.RecursivePath

		output, err := base.Run(action, true, base.LocalPath, nil, "go", "fmt", targetPath)
		if len(output) > 0 {
			action.Log(output)
		}
		if err != nil {
			action.Fatal(err)
		}
	},
})

var _ = goyek.Define(goyek.Task{
	Name:  "fmt-list",
	Usage: "Lists files that must be formated if any.",
	Action: func(action *goyek.A) {
		output, err := base.Run(action, true, base.LocalPath, nil, "gofmt", "-l", "-s", getPath(targetPath))
		if len(output) > 0 {
			action.Log(output)
		}
		if err != nil {
			action.Fatal(err)
		}
	},
})

var lint = goyek.Define(goyek.Task{
	Name:  "lint",
	Usage: `Runs "golangci-lint" on projects source code.`,
	Action: func(action *goyek.A) {
		internalPath := filepath.Join(base.LocalPath, "internal") + base.RecursivePath

		output, err := base.Run(
			action,
			true,
			base.LocalPath,
			nil,
			"go", "tool",
			fmt.Sprintf("-modfile=%s", path.Join(base.BuildPath, "golangci-lint", base.GoMod)),
			"golangci-lint",
			"run",
			internalPath,
		)
		if len(output) > 0 {
			action.Log(output)
		}
		if err != nil {
			action.Fatal(err)
		}
	},
})

var test = goyek.Define(goyek.Task{
	Name:  "test",
	Usage: "Runs unit tests.",
	Action: func(action *goyek.A) {
		internalPath := filepath.Join(base.LocalPath, "internal") + base.RecursivePath

		output, err := base.RunGoTool(
			action,
			true,
			base.LocalPath,
			nil,
			"gotestsum",
			"-f", "testname",
			"--",
			internalPath,
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
	Name:  "retest",
	Usage: `Runs "lint", "fmt" and "test" tasks.`,
	Deps: []*goyek.DefinedTask{
		lint,
		format,
		test,
	},
})
