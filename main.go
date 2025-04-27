package main

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	. "github.com/blbrdv/ezstore/internal/locale"
	"github.com/blbrdv/ezstore/internal/msstore"
	"github.com/blbrdv/ezstore/internal/windows"
	. "github.com/blbrdv/ezstore/internal/writer"
	"github.com/pterm/pterm"
	"github.com/urfave/cli/v3"
	"io"
	"os"
	"path/filepath"
	"runtime"
)

//go:embed dist/README.txt
var help string

func main() {
	defer func() {
		_ = recover()
		os.Exit(1)
	}()

	app := &cli.Command{
		Name:                  "ezstore",
		EnableShellCompletion: true,
		Writer:                DebugWriter{},
		ErrWriter:             ErrorWriter{},
		Commands: []*cli.Command{
			{
				Name:   "install",
				Action: InstallFunc,
				Arguments: []cli.Argument{
					&cli.StringArg{
						Name:  "id",
						Value: "",
					},
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "version",
						Aliases: []string{"v"},
						Value:   "latest",
						Validator: func(s string) error {
							if s == "" {
								return errors.New("value must be not empty")
							}
							return nil
						},
						ValidateDefaults: false,
					},
					&cli.StringFlag{
						Name:    "locale",
						Aliases: []string{"l"},
						Value:   "",
					},
					&cli.BoolFlag{
						Name:    "debug",
						Aliases: []string{"d"},
						Value:   false,
					},
				},
			},
		},
	}

	cli.HelpPrinter = func(_ io.Writer, _ string, _ interface{}) {
		fmt.Print(help)
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		pterm.Fatal.Println(err.Error())
	}
}

func InstallFunc(_ context.Context, cmd *cli.Command) error {
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

	locale, err := NewLocale(cmd.String("locale"))
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
