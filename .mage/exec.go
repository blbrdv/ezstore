//go:build mage

package main

import (
	"bytes"
	"github.com/magefile/mage/sh"
)

func toolV(name string, params ...string) error {
	return runTool(true, name, params...)
}

func tool(name string, params ...string) error {
	return runTool(false, name, params...)
}

func runTool(v bool, name string, params ...string) error {
	goparams := []string{"tool", modfile}
	goparams = append(goparams, name)
	goparams = append(goparams, params...)

	if v {
		return sh.RunV("go", goparams...)
	} else {
		return run("go", goparams...)
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
