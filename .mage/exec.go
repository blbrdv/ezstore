//go:build mage

package main

import (
	"bytes"
	"github.com/magefile/mage/sh"
)

func toolV(name string, params ...string) error {
	return runTool(true, name, modfile, params...)
}

func tool(name string, params ...string) error {
	return runTool(false, name, modfile, params...)
}

func runTool(v bool, name string, modfile string, params ...string) error {
	goParams := []string{"tool", modfile}
	goParams = append(goParams, name)
	goParams = append(goParams, params...)

	if v {
		return sh.RunV("go", goParams...)
	} else {
		return run("go", goParams...)
	}
}

func run(cmd string, args ...string) error {
	buff := bytes.NewBufferString("")
	_, err := sh.Exec(nil, buff, buff, cmd, args...)
	out := buff.String()
	if err != nil {
		println(out)
		return err
	}

	return nil
}
