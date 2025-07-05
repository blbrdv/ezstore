//go:build mage

package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
)

func build(target string) error {
	var err error

	exists := true
	info, err := os.Stat(target)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			exists = false
		} else {
			return err
		}
	}
	if exists && !info.IsDir() {
		return fmt.Errorf("target must be dir")
	}

	absPath, err := filepath.Abs(target)
	if err != nil {
		return err
	}

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
	err = tool("go-winres", "make", "--arch", "amd64,386,arm64,arm", "--in", "./winres.json", "--product-version", productVersion, "--file-version", fileVersion)
	if err != nil {
		return err
	}

	winresFiles, err := getWinresFiles()
	if err != nil {
		return err
	}

	err = move("./cmd", winresFiles...)
	if err != nil {
		return err
	}

	println("Compiling exe")
	err = run("go", "build", fmt.Sprintf("-ldflags=-X main.version=%s", productVersion), "-o", path.Join(absPath, "bin/ezstore.exe"), "./cmd")
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

	err = cp(absPath, files...)
	if err != nil {
		return err
	}

	return nil
}
