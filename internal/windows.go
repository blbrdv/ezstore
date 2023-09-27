package windows

import (
	"os"
	"os/exec"
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
