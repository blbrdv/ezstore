package base

import (
	"bytes"
	"fmt"
	"github.com/goyek/goyek/v2"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"unicode"
)

var envFormat = getEnvFormat()

func getEnvFormat() string {
	if runtime.GOOS == "windows" {
		return "> $env:%s = '%s'"
	} else {
		return "> export %s='%s'"
	}
}

func fold(input string) string {
	var filtered []byte

	for i := 0; i < len(input); i++ {
		if input[i] <= unicode.MaxASCII {
			filtered = append(filtered, input[i])
		}
	}

	return string(filtered)
}

// Run execute command
func Run(
	action *goyek.A,
	verbose bool,
	workDir string,
	env map[string]string,
	name string,
	args ...string,
) (output string, err error) {
	workDir = NormalizePath(workDir)

	var normalizedArgs []string
	for _, arg := range args {
		normalizedArgs = append(normalizedArgs, NormalizePath(arg))
	}

	if verbose {
		for key, value := range env {
			action.Logf(envFormat, key, value)
		}

		action.Logf("> %s %s", name, PrettifyPath(strings.Join(args, " ")))
	}

	out := bytes.NewBufferString("")
	cmd := exec.CommandContext(action.Context(), name, args...)
	cmd.Dir = workDir
	cmd.Stdin = os.Stdin
	cmd.Stdout = out
	cmd.Stderr = out
	cmd.Env = os.Environ()

	for key, value := range env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
	}

	err = cmd.Run()
	output = fold(out.String())
	output = strings.TrimSuffix(output, NewLine)

	return output, err
}
