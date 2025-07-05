//go:build mage

package main

import (
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/magefile/mage/sh"
	"path"
	"regexp"
	"strings"
)

var (
	winresFiles = []string{"rsrc_windows_386.syso", "rsrc_windows_amd64.syso"}
	goSRC       = `.\...`
	goMod       = `go.mod`
)

func removeWinresFiles() error {
	for _, file := range winresFiles {
		err := remove(path.Join("./cmd", file))
		if err != nil {
			return err
		}
	}

	return nil
}

// Clean prepares environment for clear run.
// Removes "output" and "release" directories.
// Removes winres files.
// Useless for CI runs, make sense only for local development.
func Clean() error {
	err := remove("./output", "./release")
	if err != nil {
		return err
	}

	err = removeWinresFiles()
	if err != nil {
		return err
	}

	return nil
}

// Format prettifies go source code.
// Runs "go fmt" on all go files.
//
//goland:noinspection GoUnusedExportedFunction
func Format() error {
	printf(`Formatting code in "%s"`, goSRC)
	return run("go", "fmt", goSRC)
}

// Sec run checks for known security vulnerabilities and license violations.
// Uses osv-scanner.
//
//goland:noinspection GoUnusedExportedFunction
func Sec() error {
	var err error
	lockfile := fmt.Sprintf("--lockfile=%s", goMod)

	println("Scanning dependencies for license violations")
	err = licensesScan(lockfile)
	if err != nil {
		return err
	}

	println("Scanning code for vulnerabilities")
	err = tool("osv-scanner", "scan", "source", lockfile)
	if err != nil {
		return err
	}

	return nil
}

// Lint run golangci-lint.
func Lint() error {
	println("Run linters")
	return runTool(true, "golangci-lint", `-modfile=.mage\golangci-lint\go.mod`, "run", `.\internal\...`)
}

// Check run gofmt and check if formatting needed.
func Check() error {
	println("Checking code format")
	out, err := sh.Output("gofmt", "-l", "-s", ".")
	if err != nil {
		return err
	}
	if out == "" {
		return nil
	}

	println(out)
	return fmt.Errorf("files must be formatted")
}

// Test run unit tests.
func Test() error {
	println("Running unit tests")
	return toolV("gotestsum", "-f", "testname", "--", `.\internal\...`)
}

var depRegexp = regexp.MustCompile(`^(\S+) v(\S+) v(\S+)$`)

// Deps check for dependencies updates.
// 1. Only direct dependencies.
// 2. At least minor update.
//
//goland:noinspection GoUnusedExportedFunction
func Deps() error {
	// get list of dependencies that
	// 1. need an update
	// 2. not indirect
	// in format "<name> <current version> <new version>"
	list, err := goList("{{if not (or .Indirect .Main)}}{{with .Update}}{{$.Path}} {{$.Version}} {{.Version}}{{end}}{{end}}")
	if err != nil {
		return err
	}

	var depsToUpdate []string
	for _, dep := range list {
		match := depRegexp.FindStringSubmatch(dep)
		if len(match) != 4 {
			continue
		}

		currentVersion, err := semver.NewVersion(match[2])
		if err != nil {
			return err
		}
		newVersion, err := semver.NewVersion(match[3])
		if err != nil {
			return err
		}

		if newVersion.Major() > currentVersion.Major() ||
			(newVersion.Major() == currentVersion.Major() && newVersion.Minor() > currentVersion.Minor()) {
			depsToUpdate = append(depsToUpdate, dep)
		}
	}
	if len(depsToUpdate) > 0 {
		println(strings.Join(depsToUpdate, "\n"))
		return fmt.Errorf("project dependencies must be updated")
	}

	print("No need for update!")
	return nil
}

// Build builds project.
// Compile exe and put it with necessary files into "output" dir.
func Build() error {
	return build("./output")
}

// Pack prepare project for release.
// Compiles installer and archive exe with necessary files for portable distribution.
// Put all of it to "release" dir.
func Pack() error {
	var err error

	productVersion, err := getProductVersion()
	if err != nil {
		return err
	}

	fileVersion, err := getFileVersion(productVersion)
	if err != nil {
		return err
	}

	println("Archiving files")
	err = run("7z", "a", "-bso0", "-bd", "-sse", "./release/ezstore-portable.7z", "./output/*")
	if err != nil {
		return err
	}

	println("Compiling installer")
	err = run("iscc", "/Q", "setup.iss", fmt.Sprintf("/DPV=%s", productVersion), fmt.Sprintf("/DFV=%s", fileVersion))
	if err != nil {
		return err
	}

	return nil
}
