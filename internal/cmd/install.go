package cmd

import (
	"errors"
	"fmt"
	"github.com/blbrdv/ezstore/internal/log"
	"github.com/blbrdv/ezstore/internal/ms"
	"github.com/blbrdv/ezstore/internal/ms/store"
	"github.com/blbrdv/ezstore/internal/ms/windows"
	"github.com/blbrdv/ezstore/internal/utils"
	"github.com/urfave/cli/v3"
	"golang.org/x/net/context"
)

// Install download package with its dependencies from MS Store by id, version, locale and architecture,
// and then install it all.
func Install(_ context.Context, cmd *cli.Command) error {
	if windows.Version == windows.Unsupported {
		return fmt.Errorf("unsupportded Windows version detected")
	}

	shell, err := windows.NewPowerShell()
	if err != nil {
		return err
	}
	defer func(shell *windows.Powershell) {
		shellErr := shell.Exit()
		if shellErr != nil {
			panic(shellErr)
		}
	}(shell)

	verbosity, err := log.NewLevel(cmd.String("verbosity"))
	if err != nil {
		return err
	}
	log.Level = verbosity

	id := cmd.StringArg("id")
	if id == "" {
		return errors.New("id must be set")
	}

	versionStr := cmd.String("ver")
	var version *ms.Version
	if versionStr != "latest" {
		version, err = ms.NewVersion(versionStr)
		if err != nil {
			return err
		}
	}

	var locale *ms.Locale
	if cmd.String("locale") == "" {
		locale = windows.GetLocale(shell)
	} else {
		locale, err = ms.NewLocale(cmd.String("locale"))
		if err != nil {
			return err
		}
	}

	log.Tracef("Id: %s", id)
	log.Tracef("Version: %s", version)
	log.Tracef("Locale: %s", locale)

	tmpPath := utils.Join(windows.TempDir, id)

	log.Debugf("Trace file: %s", log.TraceFile)
	log.Debugf("Temp dir: %s", tmpPath)

	defer windows.Remove(tmpPath)

	windows.Remove(tmpPath)
	windows.MkDir(tmpPath)

	file, err := store.Download(id, version, locale, tmpPath)
	if err != nil {
		return err
	}

	err = windows.Install(file, shell)
	if err != nil {
		return err
	}

	log.Success("Done!")

	return nil
}
