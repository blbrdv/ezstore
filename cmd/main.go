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

var version = "undefined"

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
			os.Exit(1)
		}
	}()

	app := &cli.Command{
		Name:                  "ezstore",
		EnableShellCompletion: true,
		Writer:                cmd.DebugWriter{},
		ErrWriter:             cmd.ErrorWriter{},
		Version:               version,
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
						Name:  "ver",
						Value: "latest",
						Validator: func(s string) error {
							if s == "" {
								return errors.New("value must be not empty")
							}
							return nil
						},
						ValidateDefaults: false,
					},
					&cli.StringFlag{
						Name:  "locale",
						Value: "",
					},
				},
			},
		},
	}

	cli.HelpPrinter = func(_ io.Writer, _ string, _ interface{}) {
		fmt.Print(help)
	}

	cli.VersionPrinter = func(cmd *cli.Command) {
		fmt.Printf("%s v%s\n", cmd.Root().Name, cmd.Root().Version)
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Errorf("Finished with error: %s", err.Error())
		os.Exit(1)
	}
}
