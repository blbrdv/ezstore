package main

import (
	"errors"
	"fmt"
	"os"
	"runtime"

	windows "github.com/blbrdv/ezstore/internal"
	"github.com/blbrdv/ezstore/internal/msstore"
	"github.com/pterm/pterm"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "ezstore",
		Usage: "Easy install apps from MS Store",
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
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		pterm.Fatal.Println(err)
	}
}

func InstallFunc(ctx *cli.Context) error {
	id := ctx.Args().Get(0)
	version := ctx.String("version")
	var arch string
	switch goarch := runtime.GOARCH; goarch {
	case "amd64":
		arch = "x64"
	case "amd64p32":
		arch = "x86"
	case "arm":
		arch = "arm"
	case "arm64":
		arch = "arm64"
	default:
		return fmt.Errorf("%s architecture not supported", goarch)
	}
	locale, err := windows.GetLocale()

	if err != nil {
		return err
	}

	if id == "" || version == "" {
		return errors.New("id and version must be set")
	}
	if ctx.String("locale") != "" {
		locale = ctx.String("locale")
	}

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
