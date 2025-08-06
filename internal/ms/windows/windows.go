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
func Install(file *ms.BundleFileInfo, shell *Powershell) error {
	var depsStr string
	deps := file.Dependencies()
	length := len(deps)

	if length == 0 {
		depsStr = ""
	} else {
		var sb strings.Builder

		_, _ = fmt.Fprint(&sb, " -DependencyPath ")

		last := length - 1
		for i := 0; i < length; i++ {
			dep := deps[i]
			_, _ = fmt.Fprintf(&sb, `"%s"`, dep.Path)
			if i != last {
				_, _ = fmt.Fprint(&sb, ", ")
			}
		}
		depsStr = sb.String()
	}

	addPkgCmd := fmt.Sprintf("Add-AppxPackage -Path %s%s", file.Path, depsStr)
	log.Tracef("Powershell: %s", addPkgCmd)
	result, err := shell.Exec(addPkgCmd)
	if err != nil {
		if result != "" {
			result = fmt.Sprintf("%s%s", utils.WindowsNewline, result)
		}
		return fmt.Errorf("can not install app %s: console command error: %s%s", file.Name, err.Error(), result)
	}

	log.Infof("Package %s %s installed.", file.Name, file.Version.String())

	return nil
}

var defaultLocale = ms.Locale{Language: "en", Country: "US"}

// GetLocale returns current locale set in hosted OS.
// If error occurred or returned value is empty, returns default locale.
func GetLocale(shell *Powershell) *ms.Locale {
	result, err := shell.Exec("Get-Culture | select -exp Name")
	if err != nil {
		return &defaultLocale
	}

	localeStr := strings.TrimSpace(result)
	localeStr = strings.Trim(localeStr, utils.WindowsNewline)

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
