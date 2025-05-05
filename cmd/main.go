package main

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/blbrdv/ezstore/internal/cmd"
	"github.com/blbrdv/ezstore/internal/log"
	"github.com/urfave/cli/v3"
	"github.com/ztrue/tracerr"
	"io"
	"os"
	"strings"
)

//go:embed README.txt
var help string

func main() {
	defer func() {
		rec := recover()
		if rec != nil {
			switch err := rec.(type) {
			case error:
				log.Errorf("Panic: %s", err.Error())

				terr := tracerr.Wrap(err)
				terr = tracerr.CustomError(err, terr.StackTrace()[2:])
				stacktrace := strings.Split(tracerr.Sprint(terr), "\n")[1:]

				fmt.Println(strings.Join(stacktrace, "\n"))
			case string:
				log.Errorf("Panic: %s", err)
			default:
				log.Error("Panic: unknown error")
			}
		}
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
						Name:  "verbosity",
						Value: "n",
					},
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
		log.Error(err.Error())
		os.Exit(1)
	}
}
