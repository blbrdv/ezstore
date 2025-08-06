package log

import (
	"fmt"
	"github.com/blbrdv/ezstore/internal/utils"
	"github.com/gookit/color"
	"os"
	"path/filepath"
	"time"
)

//goland:noinspection GoNameStartsWithPackageName
type LogLevel int

const (
	// Quiet - no output at all
	Quiet LogLevel = iota + 1
	// Minimal - only SUCCESS and ERROR logs
	Minimal
	// Normal - same as Minimal plus INFO and WARNING logs
	Normal
	// Detailed - same as Normal plus DEBUG logs and tracing net errors to log file
	Detailed
)

// Level contains current app logging level.
var Level = Normal

// TraceFile contains path to log file for http package. Can be not exists.
var TraceFile = utils.Join(
	getRoamingDir(),
	fmt.Sprintf("%s.log", time.Now().Format("060102150405")),
)

var levels = map[string]LogLevel{
	"q": Quiet,
	"m": Minimal,
	"n": Normal,
	"d": Detailed,
}

var black = color.Style{color.OpBold}
var gray = color.Style{color.FgGray, color.OpBold}
var blue = color.Style{color.FgBlue, color.OpBold}
var green = color.Style{color.FgGreen, color.OpBold}
var yellow = color.Style{color.FgYellow, color.OpBold}
var red = color.Style{color.FgRed, color.OpBold}

// NewLevel returns log level based on its string representation. Returns error if input string format is invalid.
func NewLevel(input string) (LogLevel, error) {
	result := levels[input]
	if result == 0 {
		return 0, fmt.Errorf("%s is invalid log level", input)
	}
	return result, nil
}

// Trace print input text to os.Stdout with "[TRC]" mark and appended new line when log level is Detailed.
func Trace(value string) {
	if Level == Detailed {
		utils.Fprintln(os.Stdout, black.Render("[TRC]"), value)
	}
}

// Tracef format and print input text to os.Stdout with "[TRC]" mark and appended new line when log level is Detailed.
func Tracef(format string, values ...any) {
	Trace(fmt.Sprintf(format, values...))
}

// Debug print input text to os.Stdout with "[DEB]" mark and appended new line when log level is Detailed.
func Debug(value string) {
	if Level == Detailed {
		utils.Fprintln(os.Stdout, gray.Render("[DEB]"), value)
	}
}

// Debugf format and print input text to os.Stdout with "[DEB]" mark and appended new line when log level is Detailed.
func Debugf(format string, values ...any) {
	Debug(fmt.Sprintf(format, values...))
}

// Info print input text to os.Stdout with "[INF]" mark and appended new line when log level is Normal or above.
func Info(value string) {
	if Level >= Normal {
		utils.Fprintln(os.Stdout, blue.Render("[INF]"), value)
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
		utils.Fprintln(os.Stderr, yellow.Render("[WRN]"), value)
	}
}

// Warningf format and print input text to os.Stderr with "[WRN]" mark and appended new line when log level is Normal
// or above.
func Warningf(format string, values ...any) {
	Warning(fmt.Sprintf(format, values...))
}

// Success print input text to os.Stdout with "[SCC]" mark and appended new line when log level is Minimal or above.
func Success(value string) {
	if Level >= Minimal {
		utils.Fprintln(os.Stdout, green.Render("[SCC]"), value)
	}
}

// Error print input text to os.Stderr with "[ERR]" mark and appended new line when log level is Minimal or above.
func Error(value string) {
	if Level >= Minimal {
		utils.Fprintln(os.Stderr, red.Render("[ERR]"), value)
	}
}

// Errorf format and print input text to os.Stderr with "[ERR]" mark and appended new line when log level is Minimal or
// above.
func Errorf(format string, values ...any) {
	Error(fmt.Sprintf(format, values...))
}

func getRoamingDir() string {
	roamingPath, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}

	appPath := filepath.Join(roamingPath, "ezstore")
	err = os.MkdirAll(appPath, 0660)
	if err != nil {
		panic(err.Error())
	}

	return appPath
}
