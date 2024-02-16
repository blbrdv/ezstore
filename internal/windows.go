package windows

import (
	"fmt"
	"github.com/blbrdv/ezstore/internal/types"
	"github.com/pterm/pterm"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func Install(fullPath string) error {
	regex := regexp.MustCompile(`^([0-9a-zA-Z.-]+)[-_]v?([\d\.]+)\.`)
	arr := strings.Split(fullPath, "\\")
	regexData := regex.FindStringSubmatch(arr[len(arr)-1])
	name := regexData[1]
	version, err := types.New(strings.TrimSuffix(regexData[2], "."))

	if err != nil {
		return err
	}

	cmd1 := exec.Command("powershell", "-Command",
		fmt.Sprintf(
			"Get-AppxPackage -Name %s* | Sort-Object -Property Version | Select-Object -ExpandProperty Version -Last 1",
			name))
	vBytes, err := cmd1.CombinedOutput()

	if err != nil {
		return err
	}

	vString := strings.Trim(string(vBytes), "\n\r")

	var latestVersion *types.Version

	if vString == "" {
		latestVersion, _ = types.New("0")
	} else {
		latestVersion, err = types.New(vString)

		if err != nil {
			return err
		}
	}

	if latestVersion.LessThan(version) {
		cmd2 := exec.Command("powershell", "-NoProfile", "Add-AppxPackage", "-Path", fullPath)
		cmd2err := cmd2.Run()

		if cmd2err != nil {
			return cmd2err
		}

		pterm.Debug.Println(fmt.Sprintf("Package %s v%s installed", name, version))
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
