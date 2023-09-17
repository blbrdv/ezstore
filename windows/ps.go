package windows

import (
	"github.com/KnicKnic/go-powershell/pkg/powershell"
)

func execute(command string) {
	runspace := powershell.CreateRunspaceSimple()
	defer runspace.Close()

	results := runspace.ExecScript(command, false, nil, "OS")
	results.Close()
}
