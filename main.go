package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/blbrdv/ezstore/msstore"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "ezstore",
		Usage: "Search and install apps from MS Store",
		Commands: []*cli.Command{
			{
				Name:   "search",
				Action: SearchFunc,
			},
			{
				Name:   "download",
				Action: DownloadFunc,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "version",
						Aliases: []string{"v"},
						Value:   "latest",
						Usage:   "Product version",
					},
				},
			},
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

func SearchFunc(ctx *cli.Context) error {
	args := ctx.Args()

	if args.Len() < 1 {
		return errors.New("words for search must be provided")
	}

	for _, word := range args.Slice() {
		fmt.Printf("%s\n", word)
	}

	fmt.Print("search")

	return nil
}

func DownloadFunc(ctx *cli.Context) error {
	id := ctx.Args().Get(0)
	version := ctx.String("version")

	if id == "" || version == "" {
		return errors.New("id and version must be set")
	}

	fmt.Printf("id      = %s\n", id)
	fmt.Printf("version = %s\n", version)

	fmt.Print("download")

	return nil
}

func InstallFunc(ctx *cli.Context) error {
	id := ctx.Args().Get(0)
	version := ctx.String("version")

	if id == "" || version == "" {
		return errors.New("id and version must be set")
	}

	fmt.Printf("id      = %s\n", id)
	fmt.Printf("version = %s\n", version)

	cookie, err := msstore.GetCookie()

	if err != nil {
		log.Fatal(err)
	}

	wuid, err := msstore.GetWUID(id, "US", "en")

	if err != nil {
		log.Fatal(err)
	}

	productInfos, err := msstore.GetProducts(cookie, wuid)

	if err != nil {
		log.Fatal(err)
	}

	var result []msstore.FileLocation

	for _, info := range productInfos {
		loc, err := msstore.GetUrl(info)

		if err != nil {
			log.Fatal(err)
		}

		if !strings.HasPrefix(loc.Url, "http://dl.delivery.mp.microsoft.com") {
			result = append(result, loc)
		}
	}

	for index, res := range result {
		fmt.Printf("%d\n%s\n", index, res)
	}

	return nil
}
