package windows

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func Install(fullPath string) error {
	fmt.Printf(`Installing product in "%s"`, fullPath)
	fmt.Println("")

	fmt.Printf(`Add-AppxPackage -Path %s`, fullPath)
	fmt.Println("")

	fileInfo, err := os.Stat(fullPath)

	if err != nil {
		return err
	}

	if fileInfo.IsDir() {
		fmt.Println("c")
		for _, format := range []string{"appx", "msix", "appxbundle", "msixbundle", "eappxbundle", "emsixbundle"} {
			err := installFiles(fullPath, format)

			if err != nil {
				return err
			}
		}
	} else {
		execute(fmt.Sprintf(`Add-AppxPackage -Path %s`, fullPath))
	}

	return nil
}

func installFiles(fullPath string, format string) error {
	e := filepath.Walk(fullPath, func(fullPath string, info os.FileInfo, err error) error {
		if err == nil && path.Ext(strings.ToLower(info.Name())) == format {
			execute(fmt.Sprintf(`Add-AppxPackage -Path %s`, filepath.Join(fullPath, info.Name())))
		}
		return nil
	})

	if e != nil {
		return e
	}

	return nil
}
