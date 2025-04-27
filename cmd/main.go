package main

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/blbrdv/ezstore/internal/cmd"
	"github.com/pterm/pterm"
	"github.com/urfave/cli/v3"
	"io"
	"os"
)

//go:embed README.txt
var help string

func main() {
	defer func() {
		_ = recover()
		os.Exit(1)
	}()

	app := &cli.Command{
		Name:                  "ezstore",
		EnableShellCompletion: true,
		Writer:                cmd.DebugWriter{},
		ErrWriter:             cmd.ErrorWriter{},
		Commands: []*cli.Command{
			{
				Name:   "install",
				Action: cmd.Install,
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
