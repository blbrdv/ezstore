package windows

import (
	"fmt"
	"github.com/blbrdv/ezstore/internal/log"
	"github.com/blbrdv/ezstore/internal/ms"
	"github.com/blbrdv/ezstore/internal/utils"
	"os"
	"os/exec"
	"strings"
)

// Install package if its version higher that installed counterpart.
func Install(file ms.FileInfo) error {
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
		return fmt.Errorf("can not install app %s: console command error: %s", file.Name, err.Error())
	}

	installedVersionStr := strings.Trim(string(result), "\n\r")

	var installedVersion *ms.Version

	if installedVersionStr == "" {
		installedVersion, _ = ms.NewVersion("0")
	} else {
		installedVersion, err = ms.NewVersion(installedVersionStr)
		if err != nil {
			return fmt.Errorf("can not install app %s: can not get installed version: %s", file.Name, err.Error())
		}
	}

	if installedVersion.LessThan(file.Version) {
		//cmd = exec.Command("powershell", "-NoProfile", "Add-AppxPackage", "-Path", file.Path)
		//if err = cmd.Run(); err != nil {
		//	return fmt.Errorf("can not install app %s: console command error: %s", file.Name, err.Error())
		//}

		log.Infof("Package %s %s installed.", file.Name, file.Version.String())
	} else {
		log.Infof("Package %s %s already installed. Skipping.", file.Name, installedVersion)
	}

	return nil
}

var defaultLocale = ms.Locale{Language: "en", Country: "US"}

// GetLocale returns current locale set in hosted OS.
// If error occurred or returned value is empty, returns default locale.
func GetLocale() *ms.Locale {
	cmd := exec.Command("powershell", "Get-Culture | select -exp Name")
	cultureName, err := cmd.Output()
	if err != nil {
		return &defaultLocale
	}

	localeStr := strings.TrimSpace(string(cultureName))
	localeStr = strings.Trim(localeStr, "\r\n")

	locale, err := ms.NewLocale(localeStr)
	if err != nil {
		return &defaultLocale
	}

	return locale
}

func prepareDir(dir string) string {
	dir = utils.Join(dir, "ezstore")
	err := os.MkdirAll(dir, 0660)
	if err != nil {
		panic(err)
	}
	return dir
}

func getTempDir() string {
	dir := os.TempDir()
	if dir != "" {
		return prepareDir(dir)
	}

	dir, err := os.UserCacheDir()
	if err != nil {
		panic(err)
	}

	return prepareDir(dir)
}

var TempDir = getTempDir()
