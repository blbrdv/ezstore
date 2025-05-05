package log

import (
	"fmt"
	"github.com/blbrdv/ezstore/internal/utils"
	"github.com/gookit/color"
	"os"
	"path"
	"strings"
	"time"
)

const (
	Quiet    = 1
	Minimal  = 2
	Normal   = 3
	Detailed = 4
)

func getCurrentDir() string {
	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}

	currentDir := path.Dir(strings.Replace(exePath, "\\", "/", -1))
	return currentDir
}

var Level = Normal
var TraceFile = utils.Join(getCurrentDir(), fmt.Sprintf("%s.log", time.Now().Format("060102150405")))

var levels = map[string]int{
	"q": Quiet,
	"m": Minimal,
	"n": Normal,
	"d": Detailed,
}

var Gray = color.Style{color.FgGray, color.OpBold}
var Blue = color.Style{color.FgBlue, color.OpBold}
var Green = color.Style{color.FgGreen, color.OpBold}
var Yellow = color.Style{color.FgYellow, color.OpBold}
var Red = color.Style{color.FgRed, color.OpBold}

func NewLevel(input string) (int, error) {
	result := levels[input]
	if result == 0 {
		return 0, fmt.Errorf("%s is invalid log level", input)
	}
	return result, nil
}

func Debug(value string) {
	if Level == Detailed {
		_, _ = fmt.Fprintln(os.Stdout, Gray.Render("[DEB]"), value)
	}
}

func Debugf(format string, values ...any) {
	Debug(fmt.Sprintf(format, values...))
}

func Info(value string) {
	if Level >= Normal {
		_, _ = fmt.Fprintln(os.Stdout, Blue.Render("[INF]"), value)
	}
}

func Infof(format string, values ...any) {
	Info(fmt.Sprintf(format, values...))
}

func Warning(value string) {
	if Level >= Normal {
		_, _ = fmt.Fprintln(os.Stderr, Yellow.Render("[WRN]"), value)
	}
}

func Warningf(format string, values ...any) {
	Warning(fmt.Sprintf(format, values...))
}

func Success(value string) {
	if Level >= Minimal {
		_, _ = fmt.Fprintln(os.Stdout, Green.Render("[SCC]"), value)
	}
}

func Successf(format string, values ...any) {
	Success(fmt.Sprintf(format, values...))
}

func Error(value string) {
	if Level >= Minimal {
		_, _ = fmt.Fprintln(os.Stderr, Red.Render("[ERR]"), value)
	}
}

func Errorf(format string, values ...any) {
	Error(fmt.Sprintf(format, values...))
}
