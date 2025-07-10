package base

import (
	"github.com/goyek/goyek/v2"
	"io"
	"os"
	"path"
)

// Stolen from: https://stackoverflow.com/a/74107689
func cp(srcPath, dstPath string) (err error) {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	//goland:noinspection GoUnhandledErrorResult
	defer srcFile.Close()

	dstFile, err := os.Create(dstPath)
	if err != nil {
		return err
	}

	defer func() {
		if c := dstFile.Close(); err == nil {
			err = c
		}
	}()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

func mv(action *goyek.A, delete bool, dstPath string, srcPaths ...string) (err error) {
	var verb string
	if delete {
		verb = "move"
	} else {
		verb = "copy"
	}

	dstPath = NormalizePath(dstPath)

	for _, srcPath := range srcPaths {
		srcPath = NormalizePath(srcPath)
		fullDstPath := path.Join(dstPath, path.Base(srcPath))

		action.Logf(`> %s "%s" -> "%s"`, verb, PrettifyPath(srcPath), PrettifyPath(fullDstPath))

		err = cp(srcPath, fullDstPath)
		if err != nil {
			return err
		}

		if delete {
			err = os.RemoveAll(srcPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Copy copies files to destination folder
func Copy(action *goyek.A, dst string, src ...string) error {
	return mv(action, false, dst, src...)
}

// Move moves files to destination folder
func Move(action *goyek.A, dst string, src ...string) error {
	return mv(action, true, dst, src...)
}

// Remove removes files or directories
func Remove(action *goyek.A, paths ...string) error {
	for _, targetPath := range paths {
		targetPath = NormalizePath(targetPath)

		action.Logf(`> remove "%s"`, PrettifyPath(targetPath))

		err := os.RemoveAll(targetPath)
		if err != nil {
			return err
		}
	}

	return nil
}
