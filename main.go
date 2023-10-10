package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	windows "github.com/blbrdv/ezstore/internal"
	"github.com/blbrdv/ezstore/internal/msstore"
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
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func InstallFunc(ctx *cli.Context) error {
	id := ctx.Args().Get(0)
	version := ctx.String("version")

	if id == "" || version == "" {
		return errors.New("id and version must be set")
	}

	localPath, err := os.UserCacheDir()
	fullPath := localPath + "\\ezstore\\" + id

	if err != nil {
		return err
	}

	fmt.Printf("id      = %s\n", id)
	fmt.Printf("version = %s\n", version)

	err = os.RemoveAll(fullPath)

	if err != nil {
		return err
	}

	filePath, err := msstore.Download(id, version, fullPath)

	if err != nil {
		return err
	}

	fmt.Print("Installing product ...\n")

	err = windows.Install(filePath)
	if err != nil {
		return err
	}

	fmt.Print("Done!\n")

	return nil
}
