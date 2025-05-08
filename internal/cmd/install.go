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
	"os"
)

// Install download package with its dependencies from MS Store by id, version, locale and architecture,
// and then install it all.
func Install(_ context.Context, cmd *cli.Command) error {
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
		locale = windows.GetLocale()
	} else {
		locale, err = ms.NewLocale(cmd.String("locale"))
		if err != nil {
			return err
		}
	}

	tmpPath := utils.Join(windows.TempDir, id)

	log.Debugf("Trace file: %s", log.TraceFile)
	log.Debugf("Temp dir: %s", tmpPath)

	err = os.RemoveAll(tmpPath)
	if err != nil {
		return fmt.Errorf("can not remove dir \"%s\" and its content: %s", tmpPath, err.Error())
	}

	err = os.MkdirAll(tmpPath, 0660)
	if err != nil {
		return fmt.Errorf("can not create dir \"%s\": %s", tmpPath, err.Error())
	}

	files, err := store.Download(id, version, locale, tmpPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		err = windows.Install(file)
		if err != nil {
			return err
		}
	}

	err = os.RemoveAll(tmpPath)
	if err != nil {
		return fmt.Errorf("can not remove dir \"%s\" and its content: %s", tmpPath, err.Error())
	}

	log.Success("Done!")

	return nil
}
