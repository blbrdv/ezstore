package main

import (
	"bufio"
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/magefile/mage/sh"
	"io/fs"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	winresFiles = []string{"rsrc_windows_386.syso", "rsrc_windows_amd64.syso"}
	goSRC       = `.\...`
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

// Check run multiple checks on go source code.
// Uses various linters and gofmt.
func Check() error {
	var err error

	println("Checking code for possibilities to use Go standard library")
	err = tool("usestdlibvars", goSRC)
	if err != nil {
		return err
	}

	println("Checking code for unnecessary type conversions")
	err = tool("unconvert", "-v", goSRC)
	if err != nil {
		return err
	}

	println("Checking code for unchecked errors")
	err = tool("errcheck", "-asserts", "-blank", "-ignoretests", goSRC)
	if err != nil {
		return err
	}

	println("Checking code problems")
	err = run("go", "vet", goSRC)
	if err != nil {
		return err
	}

	println("Checking code style")
	err = tool("staticcheck", goSRC)
	if err != nil {
		return err
	}

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
	return sh.RunV("go", "test", `.\internal\...`)
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
	list, err := sh.Output("go", "list", "-m", "-u", "-f", "{{if not (or .Indirect .Main)}}{{with .Update}}{{$.Path}} {{$.Version}} {{.Version}}{{end}}{{end}}", "all")
	if err != nil {
		return err
	}

	var depsToUpdate []string
	scanner := bufio.NewScanner(strings.NewReader(list))
	for scanner.Scan() {
		dep := scanner.Text()
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
	if err = scanner.Err(); err != nil {
		return err
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
	var err error

	productVersion, err := getProductVersion()
	if err != nil {
		return err
	}

	fileVersion, err := getFileVersion(productVersion)
	if err != nil {
		return err
	}

	printf("Compiling project with version %s (%s)", productVersion, fileVersion)

	println("Embedding resources")
	err = tool("go-winres", "make", "--in", "./winres.json", "--product-version", productVersion, "--file-version", fileVersion)
	if err != nil {
		return err
	}

	err = move("./cmd", winresFiles...)
	if err != nil {
		return err
	}

	println("Compiling exe")
	err = run("go", "build", fmt.Sprintf("-ldflags=-X main.version=%s", productVersion), "-o", "./output/bin/ezstore.exe", "./cmd")
	if err != nil {
		return err
	}

	err = removeWinresFiles()
	if err != nil {
		return err
	}

	var files []string
	err = filepath.WalkDir("./cmd", func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() && filepath.Ext(path) != ".go" {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return err
	}
	if len(files) < 1 {
		return fmt.Errorf("could not fild dist files in cmd directory")
	}

	err = cp("./output", files...)
	if err != nil {
		return err
	}

	return nil
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
