package main

import (
	"bytes"
	"github.com/magefile/mage/sh"
)

func tool(name string, params ...string) error {
	goparams := []string{"tool"}
	goparams = append(goparams, name)
	goparams = append(goparams, params...)

	return run("go", goparams...)
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
