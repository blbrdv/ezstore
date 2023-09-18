package windows

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
)

func Install(fullPath string) {
	fileInfo, err := os.Stat(fullPath)
	if err != nil {
		log.Fatal(err)
	}

	if fileInfo.IsDir() {
		for _, format := range []string{"appx", "appxbundle", "msix", "msixbundle"} {
			installFiles(fullPath, format)
		}
	} else {
		execute(fmt.Sprintf(`Add-AppxPackage -Path "%s"`, fullPath))
	}
}

func installFiles(fullPath string, format string) {
	e := filepath.Walk(fullPath, func(fullPath string, info os.FileInfo, err error) error {
		if err == nil && path.Ext(info.Name()) == format {
			execute(fmt.Sprintf(`Add-AppxPackage -Path "%s"`, filepath.Join(fullPath, info.Name())))
		}
		return nil
	})

	if e != nil {
		log.Fatal(e)
	}
}
