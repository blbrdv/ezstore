package cmd

import (
	"errors"
	"fmt"
	"github.com/blbrdv/ezstore/internal/log"
	"github.com/blbrdv/ezstore/internal/ms"
	"github.com/blbrdv/ezstore/internal/ms/store"
	"github.com/blbrdv/ezstore/internal/ms/windows"
	"github.com/pterm/pterm"
	"github.com/urfave/cli/v3"
	"golang.org/x/net/context"
	"os"
	"path/filepath"
	"runtime"
)

// Install download package with its dependencies from MS Store by id and version,
// and then install it with its dependencies.
func Install(_ context.Context, cmd *cli.Command) error {
	var err error
	var arch ms.Architecture

	switch goarch := runtime.GOARCH; goarch {
	case "amd64":
		arch, err = ms.NewArchitecture("x64")
	case "amd64p32":
		arch, err = ms.NewArchitecture("x86")
	case "arm", "arm64":
		arch, err = ms.NewArchitecture(goarch)
	default:
		err = fmt.Errorf("%s architecture not supported", arch)
	}
	if err != nil {
		return err
	}

	id := cmd.StringArg("id")
	if id == "" {
		return errors.New("id must be set")
	}

	versionStr := cmd.String("version")
	var version *ms.Version
	if versionStr != "latest" {
		version, err = ms.NewVersion(versionStr)
		if err != nil {
			return err
		}
	}

	var locale *ms.Locale
	if cmd.String("locale") == "" {
		locale = windows.GetLocale()
	} else {
		locale, err = ms.NewLocale(cmd.String("locale"))
		if err != nil {
			return err
		}
	}

	if log.IsTraceLevel() {
		pterm.EnableDebugMessages()
		pterm.Debug.Printfln("Trace file: %s", *log.GetLogFileName())
	} else {
		if cmd.Bool("debug") {
			pterm.EnableDebugMessages()
		}
	}

	pterm.Debug.Println(fmt.Sprintf("id           = %s", id))
	pterm.Debug.Println(fmt.Sprintf("version      = %s", versionStr))
	pterm.Debug.Println(fmt.Sprintf("locale       = %s", locale))
	pterm.Debug.Println(fmt.Sprintf("architecture = %s", arch))

	cache, _ := os.UserCacheDir()
	tmpPath := filepath.Join(cache, "ezstore", id)
	removeDir(tmpPath)

	err = os.MkdirAll(tmpPath, 0666)
	if err != nil {
		panic(err)
	}

	//defer removeDir(tmpPath)

	files, err := store.Download(id, version, arch, locale, tmpPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		err = windows.Install(file)
		if err != nil {
			return err
		}
	}

	pterm.Success.Println("Done!")

	return nil
}

func removeDir(path string) {
	err := os.RemoveAll(path)
	if err != nil {
		panic(err)
	}
}
