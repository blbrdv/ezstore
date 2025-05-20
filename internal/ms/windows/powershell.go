package windows

import (
	"fmt"
	"github.com/hnakamur/go-powershell"
)

type Powershell struct {
	*powershell.Shell
}

func (p *Powershell) Execf(format string, input ...any) (string, error) {
	return p.Shell.Exec(fmt.Sprintf(format, input...))
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
