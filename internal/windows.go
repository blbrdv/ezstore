package windows

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func Install(fullPath string) error {
	cmd := exec.Command("powershell", "-NoProfile", "Add-AppxPackage", "-Path", fullPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	if err != nil {
		return err
	}

	return nil
}

// https://stackoverflow.com/a/51831590
func GetLocale() (string, error) {
	envlang, ok := os.LookupEnv("LANG")
	if ok {
		return strings.Split(envlang, ".")[0], nil
	}

	cmd := exec.Command("powershell", "Get-Culture | select -exp Name")
	output, err := cmd.Output()
	if err == nil {
		return strings.Trim(string(output), "\r\n"), nil
	}

	return "", fmt.Errorf("cannot determine locale")
}
