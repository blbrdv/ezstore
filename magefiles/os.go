package main

import (
	"github.com/magefile/mage/sh"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

func normalizePath(path string) string {
	return strings.Replace(path, `\`, "/", -1)
}

func ppPath(path string) string {
	if runtime.GOOS == "windows" {
		return strings.Replace(path, "/", `\`, -1)
	} else {
		return strings.Replace(path, `\`, "/", -1)
	}
}

func cp(dst string, srcs ...string) error {
	return mv(false, dst, srcs...)
}

func move(dst string, srcs ...string) error {
	return mv(true, dst, srcs...)
}

func mv(delete bool, dst string, srcs ...string) error {
	var verb string
	if delete {
		verb = "Moving"
	} else {
		verb = "Copying"
	}

	pd, err := filepath.Abs(normalizePath(dst))
	if err != nil {
		return err
	}

	for _, src := range srcs {
		nsrc := normalizePath(src)
		target := path.Join(pd, path.Base(nsrc))
		ps, err := filepath.Abs(nsrc)
		if err != nil {
			return err
		}

		printf(`%s "%s" -> "%s"`, verb, ppPath(ps), ppPath(target))

		err = sh.Copy(target, ps)
		if err != nil {
			return err
		}

		if delete {
			err = sh.Rm(ps)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func remove(paths ...string) error {
	for _, pf := range paths {
		p, err := filepath.Abs(normalizePath(pf))
		if err != nil {
			return err
		}

		printf(`Removing "%s"`, ppPath(p))

		err = sh.Rm(p)
		if err != nil {
			return err
		}
	}

	return nil
}

func getNewLine() string {
	if runtime.GOOS == "windows" {
		return `\r\n`
	} else {
		return `\n`
	}
}
