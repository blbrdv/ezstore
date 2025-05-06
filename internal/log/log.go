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
	// Quiet - no output at all
	Quiet = 1
	// Minimal - only SUCCESS and ERROR logs
	Minimal = 2
	// Normal - same as Minimal plus INFO and WARNING logs
	Normal = 3
	// Detailed - same as Normal plus DEBUG logs and tracing net errors to log file
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

// Level contains current app logging level.
var Level = Normal

// TraceFile contains path to log file for http package. Can be not exists.
var TraceFile = utils.Join(getCurrentDir(), fmt.Sprintf("%s.log", time.Now().Format("060102150405")))

var levels = map[string]int{
	"q": Quiet,
	"m": Minimal,
	"n": Normal,
	"d": Detailed,
}

var gray = color.Style{color.FgGray, color.OpBold}
var blue = color.Style{color.FgBlue, color.OpBold}
var green = color.Style{color.FgGreen, color.OpBold}
var yellow = color.Style{color.FgYellow, color.OpBold}
var red = color.Style{color.FgRed, color.OpBold}

// NewLevel returns log level based on its string representation. Returns error if input string format is invalid.
func NewLevel(input string) (int, error) {
	result := levels[input]
	if result == 0 {
		return 0, fmt.Errorf("%s is invalid log level", input)
	}
	return result, nil
}

// Debug print input text to os.Stdout with "[DEB]" mark and appended new line when log level is Detailed.
func Debug(value string) {
	if Level == Detailed {
		_, _ = fmt.Fprintln(os.Stdout, gray.Render("[DEB]"), value)
	}
}

// Debugf format and print input text to os.Stdout with "[DEB]" mark and appended new line when log level is Detailed.
func Debugf(format string, values ...any) {
	Debug(fmt.Sprintf(format, values...))
}

// Info print input text to os.Stdout with "[INF]" mark and appended new line when log level is Normal or above.
func Info(value string) {
	if Level >= Normal {
		_, _ = fmt.Fprintln(os.Stdout, blue.Render("[INF]"), value)
	}
}

// Infof format and print input text to os.Stdout with "[INF]" mark and appended new line when log level is Normal or
// above.
func Infof(format string, values ...any) {
	Info(fmt.Sprintf(format, values...))
}

// Warning print input text to os.Stderr with "[WRN]" mark and appended new line when log level is Normal or above.
func Warning(value string) {
	if Level >= Normal {
		_, _ = fmt.Fprintln(os.Stderr, yellow.Render("[WRN]"), value)
	}
}

// Success print input text to os.Stdout with "[SCC]" mark and appended new line when log level is Minimal or above.
func Success(value string) {
	if Level >= Minimal {
		_, _ = fmt.Fprintln(os.Stdout, green.Render("[SCC]"), value)
	}
}

// Error print input text to os.Stderr with "[ERR]" mark and appended new line when log level is Minimal or above.
func Error(value string) {
	if Level >= Minimal {
		_, _ = fmt.Fprintln(os.Stderr, red.Render("[ERR]"), value)
	}
}

// Errorf format and print input text to os.Stderr with "[ERR]" mark and appended new line when log level is Minimal or
// above.
func Errorf(format string, values ...any) {
	Error(fmt.Sprintf(format, values...))
}
