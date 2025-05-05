package windows

import (
	"fmt"
	types "github.com/blbrdv/ezstore/internal"
	"github.com/pterm/pterm"
	"os/exec"
	"strings"
)

// Install package if its version higher that installed counterpart.
func Install(file types.FileInfo) error {
	cmd := exec.Command(
		"powershell",
		"-Command",
		fmt.Sprintf(
			"Get-AppxPackage -Name %s* | Sort-Object -Property Version | Select-Object -ExpandProperty Version -Last 1",
			file.Name,
		),
	)

	result, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	installedVersionStr := strings.Trim(string(result), "\n\r")

	var installedVersion *types.Version

	if installedVersionStr == "" {
		installedVersion, _ = types.NewVersion("0")
	} else {
		installedVersion, err = types.NewVersion(installedVersionStr)
		if err != nil {
			return err
		}
	}

	if installedVersion.LessThan(file.Version) {
		cmd = exec.Command("powershell", "-NoProfile", "Add-AppxPackage", "-Path", file.Path)
		if err = cmd.Run(); err != nil {
			return err
		}

		pterm.Info.Printfln("Package %s %s installed.", file.Name, file.Version.String())
	} else {
		pterm.Info.Printfln("Package %s %s already installed. Skipping.", file.Name, installedVersion)
	}

	return nil
}

var defaultLocale = types.Locale{Language: "en", Country: "US"}

// GetLocale returns current locale set in hosted OS.
// If error occurred or returned value is empty, returns default locale.
func GetLocale() *types.Locale {
	cmd := exec.Command("powershell", "Get-Culture | select -exp Name")
	cultureName, err := cmd.Output()
	if err != nil {
		return &defaultLocale
	}

	localeStr := strings.TrimSpace(string(cultureName))
	localeStr = strings.Trim(localeStr, "\r\n")

	locale, err := types.NewLocale(localeStr)
	if err != nil {
		return &defaultLocale
	}

	return locale
}
