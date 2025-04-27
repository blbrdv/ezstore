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

func Install(_ context.Context, cmd *cli.Command) error {
	var err error

	arch := runtime.GOARCH
	switch arch {
	case "amd64":
		arch = "x64"
	case "amd64p32":
		arch = "x86"
	case "arm", "arm64":
		break
	default:
		return fmt.Errorf("%s architecture not supported", arch)
	}

	id := cmd.StringArg("id")
	if id == "" {
		return errors.New("id must be set")
	}

	version := cmd.String("version")

	locale, err := types.NewLocale(cmd.String("locale"))
	if err != nil {
		locale = windows.GetLocale()
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
