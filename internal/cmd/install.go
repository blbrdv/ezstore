package cmd

import (
	"errors"
	"fmt"
	types "github.com/blbrdv/ezstore/internal"
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
	arch, err := types.NewArchitecture(runtime.GOARCH)
	if err != nil {
		return err
	}

	id := cmd.StringArg("id")
	if id == "" {
		return errors.New("id must be set")
	}

	version, err := types.NewVersion(cmd.String("version"))
	if err != nil {
		return err
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

	if cmd.Bool("debug") {
		pterm.EnableDebugMessages()
	}

	pterm.Debug.Println(fmt.Sprintf("id           = %s", id))
	pterm.Debug.Println(fmt.Sprintf("version      = %s", version))
	pterm.Debug.Println(fmt.Sprintf("locale       = %s", locale))
	pterm.Debug.Println(fmt.Sprintf("architecture = %s", arch))

	tmpPath := filepath.Join("%LocalAppData%", "ezstore", id)
	removeDir(tmpPath)

	defer removeDir(tmpPath)

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
