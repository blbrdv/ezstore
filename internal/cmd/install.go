package cmd

import (
	"errors"
	"fmt"
	types "github.com/blbrdv/ezstore/internal"
	"github.com/blbrdv/ezstore/internal/log"
	"github.com/blbrdv/ezstore/internal/msstore"
	"github.com/blbrdv/ezstore/internal/windows"
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
	var arch types.Architecture

	switch goarch := runtime.GOARCH; goarch {
	case "amd64":
		arch, err = types.NewArchitecture("x64")
	case "amd64p32":
		arch, err = types.NewArchitecture("x86")
	case "arm", "arm64":
		arch, err = types.NewArchitecture(goarch)
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
	var version *types.Version
	if versionStr != "latest" {
		version, err = types.NewVersion(versionStr)
		if err != nil {
			return err
		}
	}

	var locale *types.Locale
	if cmd.String("locale") == "" {
		locale = windows.GetLocale()
	} else {
		locale, err = types.NewLocale(cmd.String("locale"))
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

	files, err := msstore.Download(id, version, arch, locale, tmpPath)
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
