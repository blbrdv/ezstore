package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/blbrdv/ezstore/msstore"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "ezstore",
		Usage: "Search and install apps from MS Store",
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

	// if err := windows.Install(localPath + "\\ezstore\\" + id); err != nil {
	// 	return err
	// }

	command := fmt.Sprintf(`Add-AppxPackage -Path %s`, filePath)
	fmt.Printf("command = %s", command)

	_, err = exec.Command("powershell", "-NoProfile", command).CombinedOutput()
	if err != nil {
		return err
	}

	return nil
}
