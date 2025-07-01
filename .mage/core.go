//go:build mage

package main

import (
	"fmt"
)

var modfile = `-modfile=.mage\go.mod`

func println(value string) {
	_, err := fmt.Println(value)
	if err != nil {
		panic(err)
	}
}

func printf(format string, values ...any) {
	_, err := fmt.Printf(fmt.Sprintf("%s\n", format), values...)
	if err != nil {
		panic(err)
	}
}
