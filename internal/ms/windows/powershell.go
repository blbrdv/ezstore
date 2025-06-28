package windows

import (
	"fmt"
	"github.com/blbrdv/ezstore/internal/log"
	"github.com/hnakamur/go-powershell"
)

type Powershell struct {
	*powershell.Shell
}

func (p *Powershell) Execf(format string, input ...any) (string, error) {
	script := fmt.Sprintf(format, input...)
	log.Tracef("Powershell: %s", script)
	return p.Shell.Exec(script)
}

func (p *Powershell) Exit() {
	err := p.Shell.Exit()
	if err != nil {
		panic(err.Error())
	}
}

func NewPowershell(shell *powershell.Shell) *Powershell {
	return &Powershell{shell}
}

func getPowershell() *Powershell {
	shell, err := powershell.New()
	if err != nil {
		panic(err.Error())
	}
	return NewPowershell(shell)
}

var Shell = getPowershell()
