package main

import (
	"fmt"
	"github.com/magefile/mage/sh"
	"io/fs"
	"path/filepath"
	"strings"
)

func gofmt() error {
	params := []string{"-l"}
	err := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() && filepath.Dir(path) != "magefiles" && filepath.Ext(path) == ".go" {
			params = append(params, path)
		}
		return nil
	})
	if err != nil {
		return err
	}
	if len(params) == 1 {
		return fmt.Errorf("could not fild go source files")
	}

	out, err := sh.Output("gofmt", params...)
	if err != nil {
		return err
	}

	lines := strings.Count(out, "\n")
	if lines > 0 {
		return fmt.Errorf("files must be formatted:\n%s", out)
	}

	return nil
}
