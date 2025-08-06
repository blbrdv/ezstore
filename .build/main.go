package main

import (
	"flag"
	"fmt"
	"github.com/goyek/goyek/v2"
	"github.com/goyek/goyek/v2/middleware"
	"io"
	"os"
	"strings"
)

var (
	quiet      = flag.Bool("q", false, "No output to stdout except for errors.")
	buildPath  = flag.String("b", "output", "Path for build output.")
	packPath   = flag.String("p", "release", "Path for pack output.")
	targetPath = flag.String("t", "", "Target path.")
	targetArch = flag.String("arch", archsList, fmt.Sprintf("Target architecture. Allowed: %s", archsList))
	config     = flag.String("c", "release", fmt.Sprintf(
		"Sets level of optimisation for compilation. Allowed: %s",
		configList,
	))
)

func main() {
	goyek.SetLogger(goyek.FmtLogger{})

	out := goyek.Output()

	if err := os.Chdir(".."); err != nil {
		_, _ = fmt.Fprintln(out, err)
		os.Exit(1)
	}

	if os.Getenv("GITHUB_ACTIONS") == "" {
		goyek.Use(middleware.ReportStatus)
	}

	flag.CommandLine.SetOutput(out)
	flag.Usage = usage
	flag.Parse()

	if *quiet {
		goyek.Use(loggingMiddleware)
	}
	goyek.SetUsage(usage)
	goyek.Main(flag.Args())
}

func usage() {
	fmt.Println("Usage: [flags] [tasks]")
	goyek.Print()
	fmt.Println("Flags:")
	flag.PrintDefaults()
}

func loggingMiddleware(next goyek.Runner) goyek.Runner {
	return func(in goyek.Input) goyek.Result {
		originalOut := in.Output
		streamWriter := &strings.Builder{}
		in.Output = streamWriter

		result := next(in)
		output := streamWriter.String()

		if result.Status == goyek.StatusFailed {
			_, err := io.Copy(originalOut, strings.NewReader(output))
			if err != nil {
				panic(err)
			}
		}

		return result
	}
}

var _ = goyek.Define(goyek.Task{
	Name:  "help",
	Usage: `Prints help message.`,
	Action: func(action *goyek.A) {
		usage()
	},
})
