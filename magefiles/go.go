package main

import (
	"fmt"
	"github.com/magefile/mage/sh"
	"io"
	"regexp"
	"slices"
	"strings"
)

func goList(format string) ([]string, error) {
	output, err := sh.Output("go", "list", "-m", "-u", "-f", format, "all")
	if err != nil {
		return []string{}, err
	}

	return strings.Split(output, "\n"), nil
}

var allowedLicenses = []string{
	"0BSD",
	"MIT",
	"BSD-3-Clause",
	"Apache-2.0",
	"MPL-2.0",
	"MPL-2.0-no-copyleft-exception",
	"MPL-1.1",
	"LGPL-2.1",
	"LGPL-2.1+",
	"LGPL-3.0",
	"LGPL-3.0+",
	"GPL-2.0",
	"GPL-2.0+",
	"GPL-3.0",
	"GPL-3.0+",
	"AGPL-3.0-only",
}

var licenseRegexp = regexp.MustCompile(`^\|\s+(\S+)\s+\|\s+Go\s+\|\s+(\S+)`)

func fprint(w io.Writer, value string) {
	_, err := fmt.Fprint(w, value)
	if err != nil {
		panic(err.Error())
	}
}

func licensesScan(lockfile string) error {
	var err error
	var output string
	var sb strings.Builder

	lastIndex := len(allowedLicenses) - 1
	for index, license := range allowedLicenses {
		fprint(&sb, license)

		if index < lastIndex {
			fprint(&sb, ",")
		}
	}
	licenses := fmt.Sprintf(`--licenses=%s`, sb.String())

	packages, err := goList("{{if not (or .Indirect .Main)}}{{.Path}}{{end}}")
	if err != nil {
		return err
	}

	output, err = sh.Output("go", "tool", `-modfile=magefiles\go.mod`, "osv-scanner", licenses, lockfile)
	if err != nil && !(strings.Contains(output, "NO. OF PACKAGE VERSIONS")) {
		return err
	}
	lines := strings.Split(output, "\n")
	tableIndex := slices.IndexFunc(lines, func(line string) bool {
		return strings.Contains(line, "LICENSE VIOLATION")
	})
	if tableIndex >= 0 {
		startTableIndex := tableIndex + 2
		endTableIndex := slices.IndexFunc(lines[startTableIndex:], func(line string) bool {
			return strings.HasSuffix(line, "+")
		})

		if endTableIndex >= 0 {
			endTableIndex += startTableIndex
			found := false
			for _, line := range lines[startTableIndex:endTableIndex] {
				matches := licenseRegexp.FindStringSubmatch(line)
				if len(matches) != 3 {
					return fmt.Errorf("can not parse string: '%s'", line)
				}

				name := matches[2]
				if slices.Contains(packages, name) {
					printf("%s %s", matches[1], name)
					found = true
				}
			}

			if found {
				return fmt.Errorf("license violations found")
			}
		} else {
			return fmt.Errorf("can not parse osv-scanner output")
		}
	}

	return nil
}
