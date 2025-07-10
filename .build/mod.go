package main

import (
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/goyek/goyek/v2"
	"io/fs"
	"main/base"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
)

var (
	depRegexp       = regexp.MustCompile(`^(\S+) v(\S+) v(\S+)$`)
	allowedLicenses = []string{
		"0BSD",
		"MIT",
		"BSD-3-Clause",
		"Apache-2.0",
		"MPL-2.0",
		"MPL-2.0-no-copyleft-exception",
		"MPL-1.1",
		"LGPL-2.1",
		"LGPL-2.1+",
		"LGPL-3.0",
		"LGPL-3.0+",
		"GPL-2.0",
		"GPL-2.0+",
		"GPL-3.0",
		"GPL-3.0+",
		"AGPL-3.0-only",
	}
	licenseRegexp = regexp.MustCompile(`^\|\s+(\S+)\s+\|\s+Go\s+\|\s+(\S+)`)
)

func licensesScan(action *goyek.A, workDir, lockfile string) (err error) {
	var sb strings.Builder

	lastIndex := len(allowedLicenses) - 1
	for index, license := range allowedLicenses {
		base.Fprint(&sb, license)

		if index < lastIndex {
			base.Fprint(&sb, ",")
		}
	}
	licenses := fmt.Sprintf(`--licenses=%s`, sb.String())
	sb.Reset()

	packages, err := base.GetModulesList(action, workDir, "{{if not (or .Indirect .Main)}}{{.Path}}{{end}}")
	if err != nil {
		return base.Combine(strings.Join(packages, base.NewLine), err)
	}

	output, err := base.Run(
		action,
		true,
		workDir,
		nil,
		"go", "tool",
		base.ModFile,
		"osv-scanner",
		licenses,
		lockfile,
	)
	if err != nil && !(strings.Contains(output, "NO. OF PACKAGE VERSIONS")) {
		return base.Combine(output, err)
	}

	lines := strings.Split(output, base.NewLine)
	tableIndex := slices.IndexFunc(lines, func(line string) bool {
		return strings.Contains(line, "LICENSE VIOLATION")
	})
	if tableIndex >= 0 {
		startTableIndex := tableIndex + 2
		endTableIndex := slices.IndexFunc(lines[startTableIndex:], func(line string) bool {
			return strings.HasSuffix(line, "+")
		})

		if endTableIndex >= 0 {
			endTableIndex += startTableIndex
			found := false
			for _, line := range lines[startTableIndex:endTableIndex] {
				matches := licenseRegexp.FindStringSubmatch(line)
				if len(matches) != 3 {
					return fmt.Errorf("can not parse string: '%s'", line)
				}

				name := matches[2]
				if slices.Contains(packages, name) {
					base.Fprint(&sb, fmt.Sprintf("%s %s", matches[1], name))
					found = true
				}
			}

			if found {
				return fmt.Errorf("%s%s%s", sb.String(), base.NewLine, "license violations found")
			}
		} else {
			return fmt.Errorf("%s%s%s", output, base.NewLine, "can not parse osv-scanner output")
		}
	}

	return nil
}

var _ = goyek.Define(goyek.Task{
	Name:  "mod",
	Usage: `Execute "go mod tidy" everywhere mod files exists in project.`,
	Action: func(action *goyek.A) {
		var paths []string
		err := filepath.WalkDir(base.LocalPath, func(path string, d fs.DirEntry, err error) error {
			dir, file := filepath.Split(path)
			if !d.IsDir() && file == base.GoMod {
				paths = append(paths, dir)
			}
			return nil
		})
		if err != nil {
			action.Fatal(err)
		}
		if len(paths) < 1 {
			action.Fatal(fmt.Errorf(`could not fild mod files in "%s"`, base.LocalPath))
		}

		for _, path := range paths {
			if len(path) == 0 {
				path = base.LocalPath
			}
			action.Logf("> cd '%s'", base.PrettifyPath(path))
			output, err := base.Run(action, true, path, nil, "go", "mod", "tidy")
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
	Name:  "deps",
	Usage: "Prints dependency that must be bumped if any.",
	Action: func(action *goyek.A) {
		// get list of dependencies that
		// 1. need an update
		// 2. not indirect
		// in format "<name> <current version> <new version>"
		list, err := base.GetModulesList(
			action,
			base.LocalPath,
			"{{if not (or .Indirect .Main)}}{{with .Update}}{{$.Path}} {{$.Version}} {{.Version}}{{end}}{{end}}",
		)
		if err != nil {
			action.Fatal(err)
		}

		var depsToUpdate []string
		for _, dep := range list {
			match := depRegexp.FindStringSubmatch(dep)
			if len(match) != 4 {
				continue
			}

			currentVersion, err := semver.NewVersion(match[2])
			if err != nil {
				action.Fatal(err)
			}
			newVersion, err := semver.NewVersion(match[3])
			if err != nil {
				action.Fatal(err)
			}

			if newVersion.Major() > currentVersion.Major() ||
				(newVersion.Major() == currentVersion.Major() && newVersion.Minor() > currentVersion.Minor()) {
				depsToUpdate = append(depsToUpdate, dep)
			}
		}
		if len(depsToUpdate) > 0 {
			action.Fatal(base.Combine(strings.Join(depsToUpdate, base.NewLine), "Project dependencies must be updated."))
		}
		action.Log("No dependencies to update.")
	},
})

var _ = goyek.Define(goyek.Task{
	Name:  "sec",
	Usage: "Checks dependencies for license violations and known vulnerabilities.",
	Action: func(action *goyek.A) {
		lockfile := fmt.Sprintf("--lockfile=%s", base.GoMod)

		err := licensesScan(action, base.LocalPath, lockfile)
		if err != nil {
			action.Fatal(err)
		}
		action.Log("No license violations found.")

		output, err := base.RunGoTool(action, true, base.LocalPath, nil, "osv-scanner", "scan", "source", lockfile)
		if err != nil {
			action.Fatal(base.Combine(output, err))
		}
		action.Log("No vulnerabilities found.")
	},
})
