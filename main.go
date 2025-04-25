package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"runtime"

	windows "github.com/blbrdv/ezstore/internal"
	"github.com/blbrdv/ezstore/internal/msstore"
	"github.com/pterm/pterm"
	"github.com/urfave/cli/v3"
)

func main() {
	app := &cli.Command{
		Name:                  "ezstore",
		Usage:                 "Easy install apps from MS Store",
		EnableShellCompletion: true,
		Commands: []*cli.Command{
			{
				Name:   "install",
				Action: InstallFunc,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "version",
						Aliases: []string{"v"},
						Value:   "latest",
						Usage:   "Product version",
					},
					&cli.StringFlag{
						Name:    "locale",
						Aliases: []string{"l"},
						Value:   "",
						Usage:   "Product locale",
					},
					&cli.BoolFlag{
						Name:    "debug",
						Aliases: []string{"d"},
						Value:   false,
						Usage:   "debug output",
					},
				},
			},
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		pterm.Fatal.Println(err)
	}
}

func InstallFunc(_ context.Context, cmd *cli.Command) error {
	id := cmd.Args().Get(0)
	version := "latest"
	var arch string
	switch goarch := runtime.GOARCH; goarch {
	case "amd64":
		arch = "x64"
	case "amd64p32":
		arch = "x86"
	case "arm", "arm64":
		arch = goarch
	default:
		return fmt.Errorf("%s architecture not supported", goarch)
	}
	locale, err := windows.GetLocale()

	if err != nil {
		return err
	}

	if id == "" {
		return errors.New("id must be set")
	}
	if cmd.String("version") != "" {
		version = cmd.String("version")
	}
	if cmd.String("locale") != "" {
		locale = cmd.String("locale")
	}
	if cmd.Bool("debug") {
		pterm.EnableDebugMessages()
	}

	pterm.Debug.Println(fmt.Sprintf("id           = %s", id))
	pterm.Debug.Println(fmt.Sprintf("version      = %s", version))
	pterm.Debug.Println(fmt.Sprintf("locale       = %s", locale))
	pterm.Debug.Println(fmt.Sprintf("architecture = %s", arch))

	localPath, err := os.UserCacheDir()
	fullPath := localPath + "\\ezstore\\" + id

	if err != nil {
		return err
	}

	err = os.RemoveAll(fullPath)

	if err != nil {
		return err
	}

	files, err := msstore.Download(id, version, arch, locale, fullPath)

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
