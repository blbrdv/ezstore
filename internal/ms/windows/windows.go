package windows

import (
	"fmt"
	"github.com/blbrdv/ezstore/internal/log"
	"github.com/blbrdv/ezstore/internal/ms"
	"github.com/blbrdv/ezstore/internal/utils"
	"os"
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
	err := os.RemoveAll(path)
	if err != nil {
		panic(err.Error())
	}
}

func MkDir(path string) {
	err := os.MkdirAll(path, 0660)
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
	var result string
	var err error
	result, err = Shell.Execf(
		"Get-AppxPackage -Name %s* | Sort-Object -Property Version | Select-Object -ExpandProperty Version -Last 1",
		file.Name,
	)
	if err != nil {
		return fmt.Errorf("can not install app %s: console command error: %s", file.Name, err.Error())
	}

	installedVersionStr := strings.Trim(result, "\n\r")

	var installedVersion *ms.Version

	if installedVersionStr == "" {
		installedVersion, err = ms.NewVersion("0")
		if err != nil {
			return err
		}
	} else {
		installedVersion, err = ms.NewVersion(installedVersionStr)
		if err != nil {
			return fmt.Errorf("can not install app %s: can not get installed version: %s", file.Name, err.Error())
		}
	}

	if installedVersion.LessThan(file.Version) {
		result, err = Shell.Execf("Add-AppxPackage -Path %s", file.Path)
		if err != nil {
			if result != "" {
				result = fmt.Sprintf("\n%s", result)
			}
			return fmt.Errorf("can not install app %s: console command error: %s%s", file.Name, err.Error(), result)
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
	result, err := Shell.Exec("Get-Culture | select -exp Name")
	if err != nil {
		return &defaultLocale
	}

	localeStr := strings.TrimSpace(result)
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
