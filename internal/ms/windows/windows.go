package windows

import (
	"bytes"
	"fmt"
	"github.com/blbrdv/ezstore/internal/log"
	"github.com/blbrdv/ezstore/internal/ms"
	"github.com/blbrdv/ezstore/internal/utils"
	"golang.org/x/sys/windows"
	"os"
	"os/exec"
	"strings"
)

type File struct {
	*os.File
}

func (f *File) Close() {
	err := f.File.Close()
	if err != nil {
		panic(err.Error())
	}
}

func (f *File) WriteString(input string) {
	_, err := f.File.WriteString(input)
	if err != nil {
		panic(err.Error())
	}
}

func NewFile(file *os.File) *File {
	return &File{file}
}

func Remove(path string) {
	err := windows.Rmdir(path)
	if err != nil {
		panic(err.Error())
	}
}

func MkDir(path string) {
	err := windows.Mkdir(path, 0)
	if err != nil {
		panic(err.Error())
	}
}

func OpenFile(path string, flag int) *File {
	file, err := os.OpenFile(path, flag, 0660)
	if err != nil {
		panic(err.Error())
	}
	return NewFile(file)
}

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
		buf := new(bytes.Buffer)
		cmd = exec.Command("powershell", "-NoProfile", "Add-AppxPackage", "-Path", file.Path)
		cmd.Stderr = buf
		if err = cmd.Run(); err != nil {
			output := buf.String()
			if output != "" {
				output = fmt.Sprintf("\n%s", output)
			}
			return fmt.Errorf("can not install app %s: console command error: %s%s", file.Name, err.Error(), output)
		}

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

func prepareDir(elem ...string) string {
	dir := utils.Join(append(elem, "ezstore")...)
	MkDir(dir)
	return dir
}

func getTempDir() string {
	dir := os.TempDir()
	if dir != "" {
		return prepareDir(dir)
	}

	dir, err := os.UserCacheDir()
	if err != nil {
		panic(err.Error())
	}

	return prepareDir(dir, "Temp")
}

// TempDir contains path to directory for temporary files. Can be not exists.
var TempDir = getTempDir()
