package base

import (
	"github.com/goyek/goyek/v2"
	"strings"
)

// RunGoTool execute go tool using .build mod file
func RunGoTool(
	action *goyek.A,
	verbose bool,
	workDir string,
	env map[string]string,
	name string,
	params ...string,
) (string, error) {
	goParams := []string{"tool", ModFile}
	goParams = append(goParams, name)
	goParams = append(goParams, params...)

	return Run(action, verbose, workDir, env, "go", goParams...)
}

// GetModulesList returns result of "go list" command using provided format string
func GetModulesList(action *goyek.A, workDir, format string) ([]string, error) {
	output, err := Run(action, false, workDir, nil, "go", "list", "-m", "-u", "-f", format, "all")
	if err != nil {
		return []string{}, Combine(output, err)
	}

	return strings.Split(output, NewLine), nil
}
